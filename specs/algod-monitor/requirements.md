---
spec: algod-monitor.spec.md
---

## User Stories

- As an operator, I want deterministic algod health checks suitable for both
  interactive diagnosis and automated alerting.

## Acceptance Criteria

### REQ-algod-monitor-001

The checker SHALL classify connectivity, HTTP, response parsing, and unsupported-round failures as down with error context.

Acceptance Criteria
- Transport, non-success HTTP, malformed JSON, body-read, and unsupported-round fixtures produce a down result with a non-empty error.

### REQ-algod-monitor-002

The checker SHALL classify positive catchup time or excessive round lag as degraded and otherwise classify a successful status response as healthy.

Acceptance Criteria
- Positive catchup time and lag beyond the configured maximum produce degraded results.
- A successful response within the lag limit produces a healthy result.

### REQ-algod-monitor-003

Multi-node checks SHALL execute concurrently and preserve configuration order.

Acceptance Criteria
- Results retain the same order as their input node configurations even when response completion order differs.

### REQ-algod-monitor-004

The CLI SHALL support JSON configuration, documented defaults, flag overrides, text or JSON output, one-shot operation, and immediate-then-periodic watch mode.

Acceptance Criteria
- CLI tests demonstrate configuration loading, flag precedence, both output formats, one-shot execution, and an immediate watch-mode check before its interval tick.

### REQ-algod-monitor-005

One-shot operation SHALL fail when any checked node is degraded or down.

Acceptance Criteria
- A one-shot run containing any degraded or down node returns an error; an all-healthy run succeeds.

## Constraints

- Health uses the algod status endpoint and a caller-defined round-lag limit.

## Out of Scope

- Repairing nodes, sending alerts, and persisting historical health data.
