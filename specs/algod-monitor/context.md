---
spec: algod-monitor.spec.md
---

## Context

Operators need a small dependency-light probe that works against public and
token-protected algod nodes and can feed shell monitoring pipelines.

## Related Modules

- Algorand algod `/v2/status` HTTP contract.

## Design Decisions

- Keep status checks in an importable package and CLI/config orchestration in
  separate packages.
- Preserve input ordering even though checks execute concurrently.
