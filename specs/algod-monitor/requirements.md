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
- Existing fixtures prove request-creation, transport, non-success HTTP, malformed JSON, and unsupported-round failures produce a down result with a non-empty error.
- The response-body read failure branch assigns down status and a contextual error; it is source-verified because no committed fixture forces `io.ReadAll` to fail.

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
- Configuration tests prove defaults, JSON loading, and validation failures.
- Source inspection proves changed flags override configuration, both output branches exist, one-shot returns the result-derived error, and watch mode checks before entering its ticker loop; no committed CLI tests cover those branches.

### REQ-algod-monitor-005

One-shot operation SHALL fail when any checked node is degraded or down.

Acceptance Criteria
- `exitCodeFromResults` returns a down-node error before a degraded-node error and returns nil only when neither set is populated; this behavior is source-verified because no committed CLI test invokes it.

## Constraints

- Health uses the algod status endpoint and a caller-defined round-lag limit.

## Out of Scope

- Repairing nodes, sending alerts, and persisting historical health data.
