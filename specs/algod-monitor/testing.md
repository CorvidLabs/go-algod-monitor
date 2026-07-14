---
spec: algod-monitor.spec.md
---

## Test Plan

### Unit Tests

- Validate defaults and every invalid configuration category.
- Exercise healthy, catchup, lagged, unsupported, HTTP, transport, and malformed responses.
- Verify token-header omission/inclusion and ordered multi-node results.

### Integration Tests

- No committed CLI integration tests exist. Text/JSON output, one-shot exit,
  flag precedence, and immediate watch behavior are verified from `main.go`
  source plus successful formatting, vet, race-enabled tests, and build.
- No committed fixture forces a response-body read error; that branch is
  verified from `health.go` source plus successful compilation.
