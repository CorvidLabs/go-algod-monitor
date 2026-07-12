---
change: CHG-0001-adopt-specsync-5-0-1-and-trust-1-0-0-governance-for-go-algod-monitor
artifact: context
---

# Context

go-algod-monitor has mature native Go validation but no previous SpecSync threshold or
verified SDD policy. This migration records the existing algod monitor behavior,
adds project-scoped agent guidance, and composes the native gate through Trust
without changing runtime semantics or public APIs.
