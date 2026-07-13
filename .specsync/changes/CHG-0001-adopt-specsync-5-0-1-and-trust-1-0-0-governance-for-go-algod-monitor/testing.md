---
change: CHG-0001-adopt-specsync-5-0-1-and-trust-1-0-0-governance-for-go-algod-monitor
artifact: testing
---

# Testing

- REQ-algod-monitor-001: exercise transport, HTTP, body-read, JSON, and unsupported-round failures and require down status with error context.
- REQ-algod-monitor-002: exercise healthy, catchup, and excessive-lag classifications.
- REQ-algod-monitor-003: exercise concurrent multi-node checks and assert input-order results.
- REQ-algod-monitor-004: exercise configuration defaults, flag overrides, output formats, one-shot mode, and immediate watch-mode behavior.
- REQ-algod-monitor-005: exercise successful all-healthy one-shot execution and failure when any node is degraded or down.
- Run specsync strict validation at committed threshold zero.
- Confirm all four agent integrations report installed.
- Run the Fledge verify lane and Trust doctor.
- Confirm formatting, vet, race-enabled tests, and builds pass locally and in the hosted Trust lifecycle.
- Confirm the existing Go 1.24 and 1.25 matrix remains green.
