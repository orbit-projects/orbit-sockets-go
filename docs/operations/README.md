# Operations

Terminate public TLS at a trusted ingress or provide a separately managed TLS listener. Configure
the ingress to forward WebSocket upgrade headers, preserve the `Origin` header, and use idle and
drain timeouts suitable for long-lived sessions. Keep the service on a private network and set
`AllowedOrigins` explicitly for browser clients.

Size `MaxConnections` from measured memory and file descriptor budgets. Keep message bounds no larger
than application needs. Call `Shutdown` with a deadline on termination; the library closes active
sessions after stopping the HTTP listener. Clients remain responsible for reconnect policy.
