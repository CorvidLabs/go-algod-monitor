---
change: CHG-0002-run-the-preserved-go-matrix-when-governance-inputs-change
artifact: context
---

# Context

The accepted Trust rollout preserves the repository's existing Go 1.24 and 1.25 matrix, but the Go workflow's path filters only schedule for Go source, module, Makefile, or workflow-file changes. A governance-only rollout head therefore receives no Go matrix run even though Trust invokes the same native commands locally.

This correction makes the preserved matrix observable on governance-only pull requests. Hosted results remain external gates and are not claimed complete by this change.
