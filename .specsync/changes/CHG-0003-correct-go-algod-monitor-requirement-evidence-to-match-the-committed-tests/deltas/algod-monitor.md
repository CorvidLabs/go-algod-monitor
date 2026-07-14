## MODIFIED

### REQUIREMENT REQ-algod-monitor-001

The checker SHALL classify connectivity, HTTP, response parsing, and unsupported-round failures as down with error context.

Acceptance Criteria
- Existing fixtures prove request-creation, transport, non-success HTTP, malformed JSON, and unsupported-round failures produce a down result with a non-empty error.
- The response-body read failure branch assigns down status and a contextual error; it is source-verified because no committed fixture forces `io.ReadAll` to fail.

### REQUIREMENT REQ-algod-monitor-004

The CLI SHALL support JSON configuration, documented defaults, flag overrides, text or JSON output, one-shot operation, and immediate-then-periodic watch mode.

Acceptance Criteria
- Configuration tests prove defaults, JSON loading, and validation failures.
- Source inspection proves changed flags override configuration, both output branches exist, one-shot returns the result-derived error, and watch mode checks before entering its ticker loop; no committed CLI tests cover those branches.

### REQUIREMENT REQ-algod-monitor-005

One-shot operation SHALL fail when any checked node is degraded or down.

Acceptance Criteria
- `exitCodeFromResults` returns a down-node error before a degraded-node error and returns nil only when neither set is populated; this behavior is source-verified because no committed CLI test invokes it.
