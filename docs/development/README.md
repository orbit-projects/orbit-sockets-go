# Development

Requires Go 1.24 or later. Run `gofmt -w .`, `go test ./...`, `go test -race ./...`,
`go vet ./...`, and `go build ./...` before proposing a change.

Tests use loopback servers and deterministic handlers. Release validation additionally requires the
Autobahn WebSocket test suite, ingress upgrade/drain checks, supported Go-version CI, and realistic
connection load and soak measurements.
