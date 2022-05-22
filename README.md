# pgrok - Introspected tunnels to localhost

This is a fork of pgrok with the intention of it being used within the RESWARM application ecosystem.

```
pgrok -help                                                     
Usage: pgrok [OPTIONS] <local port or address>
Options:
  -authtoken string
    	Authentication token for identifying an pgrok account
  -config string
    	Path to pgrok configuration file. (default: $HOME/.pgrok)
  -hostname string
    	Request a custom hostname from the pgrok server. (HTTP only) (requires CNAME of your DNS)
  -httpauth string
    	username:password HTTP basic auth creds protecting the public tunnel endpoint
  -inspectaddr string
    	The addr for inspect requests (default "127.0.0.1:4040")
  -inspectpublic
    	Should export inspector to public access
  -log string
    	Write log messages to this file. 'stdout' and 'none' have special meanings (default "none")
  -log-level string
    	The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR (default "DEBUG")
  -proto string
    	The protocol of the traffic over the tunnel (http+https|https|tcp) (default "http+https")
  -serveraddr string
    	The addr for server
  -subdomain string
    	Request a custom subdomain from the pgrok server. (HTTP only)
  -tls
    	Use dial for tls port
  -tlsClientCrt string
    	Path to a TLS Client CRT file if server requires
  -tlsClientKey string
    	Path to a TLS Client Key file if server requires

Examples:
	pgrok 80
	pgrok -subdomain=example 8080
	pgrok -proto=tcp 22
	pgrok -hostname="example.com" -httpauth="user:password" 10.0.0.1


Advanced usage: pgrok [OPTIONS] <command> [command args] [...]
Commands:
	pgrok start [tunnel] [...]    Start tunnels by name from config file
	ngork start-all               Start all tunnels defined in config file
	pgrok list                    List tunnel names from config file
	pgrok help                    Print help
	pgrok version                 Print pgrok version

Examples:
	pgrok start www api blog pubsub
	pgrok -log=stdout -config=pgrok.yml start ssh
	pgrok start-all
	pgrok version
```
