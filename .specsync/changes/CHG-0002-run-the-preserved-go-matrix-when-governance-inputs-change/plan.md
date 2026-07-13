---
change: CHG-0002-run-the-preserved-go-matrix-when-governance-inputs-change
artifact: plan
---

# Plan

1. Obtain portable definition approval for this exact single-file scheduling correction.
2. Add the four governance paths under both `push.paths` and `pull_request.paths`.
3. Inspect the diff to prove no job or matrix change occurred.
4. Run PR #353 SpecSync strict validation without cache, agent status, Trust doctor, and full native Trust verification.
5. Only after successful local verification, obtain closing approval and accept the change.
6. Publish the verified head and treat hosted Trust, the Go matrix, and independent review as external merge gates.
