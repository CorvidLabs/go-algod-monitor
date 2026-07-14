---
change: CHG-0002-run-the-preserved-go-matrix-when-governance-inputs-change
artifact: design
---

# Design

Add exactly four path entries to each existing event filter in `.github/workflows/go.yml`:

- `.trust.toml`
- `.specsync/**`
- `fledge.toml`
- `.github/workflows/trust.yml`

Do not change job definitions, matrix versions, steps, runner selection, permissions, concurrency, or native commands. No canonical product specification changes.
