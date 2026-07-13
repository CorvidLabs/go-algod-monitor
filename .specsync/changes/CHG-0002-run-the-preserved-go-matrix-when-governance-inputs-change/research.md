---
change: CHG-0002-run-the-preserved-go-matrix-when-governance-inputs-change
artifact: research
---

# Research

- Confirmed the accepted rollout head has no hosted Go workflow run: the branch query for workflow `Go` returns an empty set.
- Inspected `.github/workflows/go.yml`: both `push.paths` and `pull_request.paths` omit all new governance inputs.
- Confirmed the workflow's jobs, steps, runner, timeouts, and Go 1.24 and 1.25 matrix are independent of this scheduling correction.
- Identified the minimal scheduling inputs shared by the rollout: `.trust.toml`, `.specsync/**`, `fledge.toml`, and `.github/workflows/trust.yml`.
