---
module: algod-monitor
version: 1
status: active
files:
  - cmd/algod-monitor/main.go
  - internal/config/config.go
  - pkg/health/health.go

db_tables: []
depends_on: []
---

# Algod-monitor

## Purpose

Provides a command-line monitor for one or more Algorand algod nodes. It loads
JSON configuration or safe public defaults, checks node status concurrently,
classifies health from connectivity, catchup state, and round lag, and emits
human-readable or JSON results for operators and automation.

## Public API

### Exported Functions

| Export | Description |
|--------|-------------|
| `DefaultConfig` | Return the two public AlgoNode defaults |
| `Load` | Load node settings from a JSON file |
| `Validate` | Validate node presence, addresses, and interval |
| `String` | Return the stable text form of a health status |
| `NewChecker` | Construct a checker with default timeout and lag values when omitted |
| `Check` | Check one algod node |
| `CheckMultiple` | Check nodes concurrently while preserving order |
| `Config` | Top-level monitor configuration |
| `Status` | Node health classification |
| `NodeStatus` | Serializable check result |
| `Checker` | Node health checker |
| `NodeConfig` | Node connection settings |
| `StatusHealthy` | Healthy status constant |
| `StatusDegraded` | Degraded status constant |
| `StatusDown` | Down status constant |

### Structs & Enums

| Type | Description |
|------|-------------|
| `Config` | Nodes, interval, maximum lag, and output settings |
| `Status` | Healthy, degraded, or down classification |
| `NodeConfig` | Address, optional API token, and display name |
| `NodeStatus` | Serializable health result including round, lag, catchup, latency, version, and error |
| `Checker` | Concurrent-safe node health checker |

### Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `Checker.Check` | `(context.Context, address, token string) NodeStatus` | Check one algod status endpoint |
| `Checker.CheckMultiple` | `(context.Context, []NodeConfig) []NodeStatus` | Check nodes concurrently while preserving input order |
| `Config.Validate` | `() error` | Reject missing nodes, empty addresses, and intervals below one second |

### Constants

| Constant | Description |
|----------|-------------|
| `StatusHealthy` | Successful and current node |
| `StatusDegraded` | Responding node that is catching up or lagging |
| `StatusDown` | Unreachable, invalid, or unsupported node |

## Invariants

1. Requests target `{address}/v2/status` and include `X-Algo-API-Token` only when a token is configured.
2. Transport, non-200, body-read, JSON, and unsupported-round failures classify the node as down and retain a useful error.
3. A responding node is degraded when catchup time is positive or time since its last round exceeds the configured maximum; otherwise it is healthy.
4. Multi-node checks run concurrently and return results in the same order as their input configurations.
5. CLI flags override file/default interval, output, and maximum-lag values.
6. One-shot execution returns an error when any node is degraded or down; watch mode runs immediately and then at the configured interval until interrupted.

## Behavioral Examples

```
Given two configured nodes where one is current and one reports catchup time
When the checker runs both nodes
Then results retain configuration order and classify them healthy and degraded
```

## Error Cases

| Error | When | Behavior |
|-------|------|----------|
| Configuration read or parse error | A requested file is missing or malformed | Return contextual error without starting checks |
| Invalid configuration | Nodes are absent, an address is empty, or interval is below one | Reject the configuration |
| Node down | Request creation, transport, status, body, JSON, or unsupported-round failure | Return a down result with error context |

## Dependencies

- Go standard library
- `github.com/spf13/cobra`

## Change Log

| Version | Date | Changes |
|---------|------|---------|
| 1 | 2026-07-12 | Initial active specification of existing monitor behavior |
