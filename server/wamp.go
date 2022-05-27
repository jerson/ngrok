package server

import (
	"context"
	"crypto/tls"
	"errors"
	"os"
	"pgrok/log"
	"time"

	"github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/wamp"
	"github.com/gammazero/nexus/v3/wamp/crsign"
)

type WampSession struct {
	client *client.Client
}

type wampLogWrapper struct {
	logger log.Logger
}

type DeviceStatus string

const (
	CONNECTED    DeviceStatus = "CONNECTED"
	DISCONNECTED DeviceStatus = "DISCONNECTED"
)

var ErrNotConnected = errors.New("not connected")

func newWampLogger(logger log.Logger) wampLogWrapper {
	return wampLogWrapper{logger: logger}
}

func (wl wampLogWrapper) Print(v ...interface{}) {
	wl.logger.Debug("%+v", v)
}

func (wl wampLogWrapper) Println(v ...interface{}) {
	wl.logger.Debug("%+v\n", v)
}

func (wl wampLogWrapper) Printf(format string, v ...interface{}) {
	wl.logger.Debug(format, v)
}

func wrapZeroLogger(logger log.Logger) wampLogWrapper {
	wrapper := newWampLogger(logger)
	return wrapper
}

type SocketConfig struct {
	PingPongTimeout   time.Duration
	ResponseTimeout   time.Duration
	ConnectionTimeout time.Duration
	SetupTestament    bool
}

func createConnectConfig() (*client.Config, error) {

	cfg := client.Config{
		Realm: "realm1",
		TlsCfg: &tls.Config{
			InsecureSkipVerify: true,
		},
		HelloDetails: wamp.Dict{
			"authid": "system",
		},
		AuthHandlers: map[string]client.AuthFunc{
			"wampcra": clientAuthFunc(),
		},
		Debug:  false,
		Logger: wrapZeroLogger(log.NewPrefixLogger("wamp", "server")),
	}

	return &cfg, nil
}

// New creates a new wamp session from a ReswarmConfig file
func NewWamp() (*WampSession, error) {
	session := &WampSession{}
	clientChannel := EstablishSocketConnection()

	client := <-clientChannel
	session.client = client

	return session, nil
}

func (wampSession *WampSession) Reconnect() {
	wampSession.Close()

	clientChannel := EstablishSocketConnection()
	client := <-clientChannel
	wampSession.client = client

}

func (wampSession *WampSession) Publish(topic Topic, args []interface{}, kwargs Dict, options Dict) error {
	if !wampSession.Connected() {
		return ErrNotConnected
	}

	return wampSession.client.Publish(string(topic), wamp.Dict(options), args, wamp.Dict(kwargs))
}

func EstablishSocketConnection() chan *client.Client {
	resChan := make(chan *client.Client)
	reswarmCrossbarURI := os.Getenv("CROSSBAR_URI")
	systemSecret := os.Getenv("SYSTEM_SECRET")

	log.Info("Attempting to establish a socket connection to %s with secret %s...", reswarmCrossbarURI, systemSecret)

	go func() {
		for {
			connectionConfig, err := createConnectConfig()
			if err != nil {
				log.Error("%s", "failed to create connect config...")
				continue
			}

			var ctx context.Context
			var cancelFunc context.CancelFunc

			ctx, cancelFunc = context.WithTimeout(context.Background(), 1250)

			var duration time.Duration
			requestStart := time.Now() // time request
			wClient, err := client.ConnectNet(ctx, reswarmCrossbarURI, *connectionConfig)
			if err != nil {
				if cancelFunc != nil {
					cancelFunc()
				}

				duration = time.Since(requestStart)

				log.Debug("Failed to establish connection: %v, retrying... (duration: %v)", err.Error(), duration.String())

				time.Sleep(time.Millisecond * 100)
				continue
			}

			if cancelFunc != nil {
				cancelFunc()
			}

			if wClient.Connected() {
				duration = time.Since(requestStart)
				log.Debug("Sucessfully established a connection (duration: %s)", duration.String())
				resChan <- wClient
				close(resChan)
				return
			}

			duration = time.Since(requestStart)
			if wClient != nil {
				wClient.Close()
			}

			log.Debug("A Session was established, but we are not connected (duration: %s)", duration.String())
		}
	}()

	return resChan
}

func (wampSession *WampSession) Connected() bool {
	if wampSession.client == nil {
		return false
	}
	return wampSession.client.Connected()
}
func (wampSession *WampSession) Done() <-chan struct{} {
	return wampSession.client.Done()
}

func (wampSession *WampSession) Subscribe(topic Topic, cb func(Result) error, options Dict) error {
	handler := func(event *wamp.Event) {
		cbEventMap := Result{
			Subscription: uint64(event.Subscription),
			Publication:  uint64(event.Publication),
			Details:      Dict(event.Details),
			Arguments:    []interface{}(event.Arguments),
			ArgumentsKw:  Dict(event.ArgumentsKw),
		}
		err := cb(cbEventMap)
		if err != nil {
			log.Error("An error occured during the subscribe result of %s", topic)
		}
	}

	return wampSession.client.Subscribe(string(topic), handler, wamp.Dict(options))
}

func (wampSession *WampSession) SubscriptionID(topic Topic) (id uint64, ok bool) {
	subID, ok := wampSession.client.SubscriptionID(string(topic))
	return uint64(subID), ok
}
func (wampSession *WampSession) RegistrationID(topic Topic) (id uint64, ok bool) {
	subID, ok := wampSession.client.RegistrationID(string(topic))
	return uint64(subID), ok
}

func (wampSession *WampSession) Call(
	ctx context.Context,
	topic Topic,
	args []interface{},
	kwargs Dict,
	options Dict,
	progCb func(Result)) (Result, error) {

	if !wampSession.Connected() {
		return Result{}, ErrNotConnected
	}

	var handler func(result *wamp.Result)
	if progCb != nil {
		handler = func(result *wamp.Result) {
			cbResultMap := Result{
				Request:     uint64(result.Request),
				Details:     Dict(result.Details),
				Arguments:   []interface{}(result.Arguments),
				ArgumentsKw: Dict(result.ArgumentsKw),
			}
			progCb(cbResultMap)
		}
	}

	result, err := wampSession.client.Call(ctx, string(topic), wamp.Dict(options), args, wamp.Dict(kwargs), handler)
	if err != nil {
		return Result{}, err
	}

	callResultMap := Result{
		Request:     uint64(result.Request),
		Details:     Dict(result.Details),
		Arguments:   []interface{}(result.Arguments),
		ArgumentsKw: Dict(result.ArgumentsKw),
	}

	return callResultMap, nil
}

func (wampSession *WampSession) GetSessionID() uint64 {
	if !wampSession.Connected() {
		return 0
	}

	return uint64(wampSession.client.ID())
}

func (wampSession *WampSession) Register(topic Topic, cb func(ctx context.Context, invocation Result) (*InvokeResult, error), options Dict) error {

	invocationHandler := func(ctx context.Context, invocation *wamp.Invocation) client.InvokeResult {
		cbInvocationMap := Result{
			Request:      uint64(invocation.Request),
			Registration: uint64(invocation.Registration),
			Details:      Dict(invocation.Details),
			Arguments:    invocation.Arguments,
			ArgumentsKw:  Dict(invocation.ArgumentsKw),
		}

		resultMap, invokeErr := cb(ctx, cbInvocationMap)
		if invokeErr != nil {
			// Global error logging for any Registered WAMP topics
			log.Error("An error occured during invocation of %s", topic)

			return client.InvokeResult{
				Err: wamp.URI("wamp.error.canceled"), // TODO: parse Error URI from error
				Args: wamp.List{
					wamp.Dict{"error": invokeErr.Error()},
				},
			}
		}

		kwargs := resultMap.ArgumentsKw

		return client.InvokeResult{Args: resultMap.Arguments, Kwargs: wamp.Dict(kwargs)}
	}

	err := wampSession.client.Register(string(topic), invocationHandler, wamp.Dict{"force_reregister": true})
	if err != nil {
		return err
	}

	return nil
}

func (wampSession *WampSession) Unregister(topic Topic) error {
	return wampSession.client.Unregister(string(topic))
}

func (wampSession *WampSession) Unsubscribe(topic Topic) error {
	return wampSession.client.Unsubscribe(string(topic))
}

func (wampSession *WampSession) Close() {
	if wampSession.client != nil {
		wampSession.client.Close() // only possible error is if it's already closed
		wampSession.client = nil
	}
}

func clientAuthFunc() func(c *wamp.Challenge) (string, wamp.Dict) {
	return func(c *wamp.Challenge) (string, wamp.Dict) {
		return crsign.RespondChallenge(os.Getenv("SYSTEM_SECRET"), c, nil), wamp.Dict{}
	}
}
