---
spec: algod-monitor.spec.md
---

## User Stories

- As an operator, I want deterministic algod health checks suitable for both
  interactive diagnosis and automated alerting.

## Acceptance Criteria

### REQ-algod-monitor-001

The checker SHALL classify connectivity, HTTP, response parsing, and
unsupported-round failures as down with error context.

### REQ-algod-monitor-002

The checker SHALL classify positive catchup time or excessive round lag as
degraded and otherwise classify a successful status response as healthy.

### REQ-algod-monitor-003

Multi-node checks SHALL execute concurrently and preserve configuration order.

### REQ-algod-monitor-004

The CLI SHALL support JSON configuration, documented defaults, flag overrides,
text or JSON output, one-shot operation, and immediate-then-periodic watch mode.

### REQ-algod-monitor-005

One-shot operation SHALL fail when any checked node is degraded or down.

## Constraints

- Health uses the algod status endpoint and a caller-defined round-lag limit.

## Out of Scope

- Repairing nodes, sending alerts, and persisting historical health data.
