# go-algod-monitor

A lightweight CLI tool for monitoring Algorand node health. Checks algod `/v2/status` endpoints, reports round lag, latency, and catchup state.

## Install

```bash
go install github.com/CorvidLabs/go-algod-monitor/cmd/algod-monitor@latest
```

## Usage

```bash
# One-shot check against default public nodes
algod-monitor --once

# Watch mode (every 30s)
algod-monitor

# Custom config file
algod-monitor -c config.json

# JSON output for scripting / monitoring pipelines
algod-monitor --once -o json
```

## Config file

```json
{
  "nodes": [
    {"address": "https://mainnet-api.algonode.cloud", "name": "algonode-mainnet"},
    {"address": "http://localhost:4001", "token": "aaaa...", "name": "local"}
  ],
  "interval_sec": 30,
  "max_lag_sec": 30,
  "output": "text"
}
```

## Build

```bash
make build    # binary in bin/
make test     # run tests with race detector
make lint     # vet + gofmt check
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0    | All nodes healthy |
| 1    | One or more nodes degraded or down |

Useful for scripted health checks and alerting pipelines.

## License

MIT
