---
change: CHG-0001-adopt-specsync-5-0-1-and-trust-1-0-0-governance-for-go-algod-monitor
artifact: research
---

# Research

- Existing CI preserves its Go 1.24 and 1.25 matrix.
- Native verification consists of formatting, vet, race-enabled tests, and builds.
- No prior SpecSync threshold exists, so contract coverage starts advisory at
  zero while strict spec validity and native verification remain blocking.
- The repository has no separate Atlas publication or provenance-signing flow.
