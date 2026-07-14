---
change: CHG-0003-correct-go-algod-monitor-requirement-evidence-to-match-the-committed-tests
artifact: testing
---

# Testing

- `REQ-algod-monitor-001`: automated evidence is `TestCheck_Down_BadURL`, `TestCheck_Down_Unreachable`, `TestCheck_Down_BadStatus`, `TestCheck_Down_InvalidJSON`, and `TestCheck_Down_StoppedAtUnsupported`; source inspection covers the unforced body-read error branch.
- `REQ-algod-monitor-002`: automated evidence is `TestCheck_Healthy`, `TestCheck_Degraded_Catchup`, and `TestCheck_Degraded_RoundLag`.
- `REQ-algod-monitor-003`: automated evidence is `TestCheckMultiple`.
- `REQ-algod-monitor-004`: automated evidence is the eight `internal/config` tests; source inspection covers CLI flag precedence, output selection, one-shot execution, and immediate watch execution because no CLI test exists.
- `REQ-algod-monitor-005`: source inspection covers `exitCodeFromResults` because no committed test invokes that function.
- The remaining health tests cover status strings, token inclusion/omission, checker defaults, and context cancellation.
- Run formatting, vet, all 22 race-enabled tests, build, strict SpecSync at 100%, and the complete Trust gate.
