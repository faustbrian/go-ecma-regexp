# Limits and security

ECMAScript backreferences and lookaround require a backtracking engine and do
not have RE2's linear-time guarantee. Treat untrusted patterns and inputs as
potentially adversarial.

`ParseOptions` bounds pattern bytes, AST depth and nodes, captures, and
character classes. `CompileOptions` bounds program instructions.
`MatchOptions` bounds input bytes and runes, VM steps, backtracks, stack and
recursion depth, logical allocations, results, replacement/split output, and
wall time. Zero means no allowance, not unlimited.

Replacement templates independently use the same input byte/code-point
allowances as subjects. Consumed token and name work is charged even when it
produces no output. Admission bounds temporary name storage before allocation;
callers migrating from subject-only budgets must allow both inputs explicitly.

Use a request-scoped context in addition to finite budgets. Cancellation,
wall-time exhaustion, and resource exhaustion are reported separately as
`context` errors, `TimeoutError`, and `LimitError`. Do not convert these errors
to a normal non-match without an explicit application policy.

All public execution entry points check cancellation before input admission
and poll during UTF-8 and UTF-16 preparation. Cancellation observed during
preparation returns no partial result and does not enter the VM. Session state
remains unchanged when cancellation is observed before an execution result.
Cancellation is cooperative, not an asynchronous interruption of each memory
operation. `MatchLimits.WallTime` begins at VM execution; use a caller deadline
to include subject and replacement preparation in the operation budget.

Execution is synchronous. The package creates no timeout goroutine, uses no
`unsafe`, and maintains no global mutable cache. Applications may add a
caller-owned bounded cache of immutable `Program` values.

`make hostile` exercises catastrophic backtracking, zero-width loops, capture
and replacement growth, nested assertions, Unicode sets, and malformed UTF-8.
`make leak` repeats budget failures under goroutine and retained-heap checks.
`make race` shares immutable programs through a caller-owned synchronized cache.

Default limits are safe starting points, not universal policy. Lower them for
interactive or multi-tenant validation and measure real workloads before
raising them.

The versioned [threat model](threat-model.md) records current boundaries,
accepted risks, and the pending replacement-work repair verification.

Three exact `G115` source dispositions cover nonnegative, token-bounded
capture-count conversions and the pinned Unicode table sentinel encoding.
They do not exclude files or other scanner rules. The library maintainer
must review them when capture counting, token consumption, or generated
Unicode inputs change; current aliases select tables 0 through 434 and encode
them as values 1 through 435 in a `uint16` field.
