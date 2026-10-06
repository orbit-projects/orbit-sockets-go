# 0001: Keep WebSocket listeners independent of the Core process host

- Status: accepted
- Date: 2026-10-06

## Context

Core's gRPC process host supervises local plugins. It is not the application's distributed RPC
system and does not define a browser-connection stream contract.

## Decision

This package is a standalone Go WebSocket service deployed behind a WebSocket-capable ingress.
Applications select it explicitly. It enforces bounded sessions and messages, validates paths and
browser origins, propagates request cancellation, and exposes deadline-bound graceful shutdown.

## Consequences

Deployments own addressing, ingress, TLS termination, and horizontal routing. Shared conformance
tests are required before declaring the Go and Python implementations interchangeable. Core's
local process-plugin transport remains separate.
