---
change: CHG-0002-run-the-preserved-go-matrix-when-governance-inputs-change
artifact: testing
---

# Testing

- Run PR #353 `specsync check --strict --force`.
- Confirm all four native agent integrations remain installed.
- Run Trust doctor and the repository's full native Trust verification.
- Inspect `git diff` to confirm only eight path-list lines changed: four under `push` and four under `pull_request`.
- Confirm the preserved Go 1.24 and 1.25 matrix and every job step are byte-for-byte unchanged.
- After publication, require hosted Trust and every Go 1.24 and 1.25 matrix leg to pass on the current head; do not record those external results early.
