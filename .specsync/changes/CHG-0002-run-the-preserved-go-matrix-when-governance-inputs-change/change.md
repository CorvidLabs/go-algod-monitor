---
id: CHG-0002-run-the-preserved-go-matrix-when-governance-inputs-change
state: accepted
type: bug_fix
base_commit: 68efb4b83bb65c1a12d97acc9815d1bc33558fe9
---

# Run the preserved Go matrix when governance inputs change

## Intent

Run the preserved Go matrix when governance inputs change

## Affected Canonical Specs

- None

## Acceptance Criteria

- Both push and pull_request path filters include .trust.toml; .specsync/**; fledge.toml; and .github/workflows/trust.yml while every Go job and matrix version remains unchanged; local strict SpecSync and native Trust verification pass before closing approval; hosted matrix results remain external

## No-spec Rationale

This correction changes only CI path-filter scheduling so the repository's preserved Go matrix validates governance-only pull-request heads; jobs, matrix versions, runtime behavior, and canonical product contracts remain unchanged.
