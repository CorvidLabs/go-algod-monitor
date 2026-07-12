---
spec: algod-monitor.spec.md
---

## Test Plan

### Unit Tests

- Validate defaults and every invalid configuration category.
- Exercise healthy, catchup, lagged, unsupported, HTTP, transport, and malformed responses.
- Verify token-header omission/inclusion and ordered multi-node results.

### Integration Tests

- Run the CLI against deterministic local HTTP fixtures and validate text/JSON
  output and one-shot exit behavior.
