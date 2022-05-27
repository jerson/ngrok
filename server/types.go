package server

import (
	"context"
)

type Dict map[string]interface{}

// Result the value that is received from a Messenger request
// E.g. Call, Register callback, Subscription callback
type Result struct {
	Arguments    []interface{}
	ArgumentsKw  Dict
	Details      Dict
	Registration uint64 // Values will be filled based on Result type, e.g. Call, Subcribe, Register...
	Request      uint64
	Subscription uint64
	Publication  uint64
}

// InvokeResult the value that is sent whenever a procedure is invoked
type InvokeResult struct {
	Arguments   []interface{}
	ArgumentsKw Dict
	Err         string
}

type Messenger interface {
	Register(topic Topic, cb func(ctx context.Context, invocation Result) (*InvokeResult, error), options Dict) error
	Publish(topic Topic, args []interface{}, kwargs Dict, options Dict) error
	Subscribe(topic Topic, cb func(Result) error, options Dict) error
	Call(ctx context.Context, topic Topic, args []interface{}, kwargs Dict, options Dict, progCb func(Result)) (Result, error)
	SubscriptionID(topic Topic) (id uint64, ok bool)
	RegistrationID(topic Topic) (id uint64, ok bool)
	Unregister(topic Topic) error
	Unsubscribe(topic Topic) error
	SetupTestament() error
	GetSessionID() uint64
	Done() <-chan struct{}
	Connected() bool
	Reconnect()
	Close()
}
