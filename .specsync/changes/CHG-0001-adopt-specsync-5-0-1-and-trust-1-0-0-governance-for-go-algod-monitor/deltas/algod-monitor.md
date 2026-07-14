## MODIFIED

### SPEC SECTION Invariants

1. Requests target `{address}/v2/status` and include `X-Algo-API-Token` only when a token is configured.
2. Transport, non-200, body-read, JSON, and unsupported-round failures classify the node as down and retain a useful error.
3. A responding node is degraded when catchup time is positive or time since its last round exceeds the configured maximum; otherwise it is healthy.
4. Multi-node checks run concurrently and return results in the same order as their input configurations.
5. CLI flags override file/default interval, output, and maximum-lag values.
6. One-shot execution returns an error when any node is degraded or down; watch mode runs immediately and then at the configured interval until interrupted.

### REQUIREMENT REQ-algod-monitor-001

The checker SHALL classify connectivity, HTTP, response parsing, and unsupported-round failures as down with error context.

Acceptance Criteria
- Transport, non-success HTTP, malformed JSON, body-read, and unsupported-round fixtures produce a down result with a non-empty error.

### REQUIREMENT REQ-algod-monitor-002

The checker SHALL classify positive catchup time or excessive round lag as degraded and otherwise classify a successful status response as healthy.

Acceptance Criteria
- Positive catchup time and lag beyond the configured maximum produce degraded results.
- A successful response within the lag limit produces a healthy result.

### REQUIREMENT REQ-algod-monitor-003

Multi-node checks SHALL execute concurrently and preserve configuration order.

Acceptance Criteria
- Results retain the same order as their input node configurations even when response completion order differs.

### REQUIREMENT REQ-algod-monitor-004

The CLI SHALL support JSON configuration, documented defaults, flag overrides, text or JSON output, one-shot operation, and immediate-then-periodic watch mode.

Acceptance Criteria
- CLI tests demonstrate configuration loading, flag precedence, both output formats, one-shot execution, and an immediate watch-mode check before its interval tick.

### REQUIREMENT REQ-algod-monitor-005

One-shot operation SHALL fail when any checked node is degraded or down.

Acceptance Criteria
- A one-shot run containing any degraded or down node returns an error; an all-healthy run succeeds.
