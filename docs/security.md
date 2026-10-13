# Limits and security

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API and indices](api.md)

ECMAScript backreferences and lookaround require a backtracking engine and do
not have RE2's linear-time guarantee. Treat untrusted patterns and inputs as
potentially adversarial.

`ParseOptions` bounds pattern bytes, AST depth and nodes, captures, and
character classes. `CompileOptions` bounds emitted program instructions,
not all compile work.
A mandatory repetition whose body emits no instructions is compiled once, not
expanded count times, including nested empty bodies. Captures, assertions and
alternative branches still emit instructions and retain their normal budgets
and semantics. Parse/Compile have no context API or independent CPU deadline;
class-set preparation and individual AST traversals are not instruction charges.
`MatchOptions` bounds input bytes and runes, VM steps, backtracks, stack and
recursion depth, logical allocations, results, replacement/split output, and
wall time. Zero means no allowance, not unlimited. These counters are selected
logical work boundaries: Steps do not charge every comparison inside an
instruction, and Allocations do not count all Go heap allocations or bytes.
Preparation, successful capture copies and output storage need their own
application capacity policy; this is not a process memory ceiling.

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

The explicit `make check-full` test operation exercises catastrophic
backtracking, zero-width loops, capture and replacement growth, nested
assertions, Unicode sets and malformed UTF-8. It retains the ten-count
goroutine/retained-heap leak assertion. Bounded `make check` selects the existing
immutable-program/caller-owned-cache race assertion once with a two-minute
timeout. See [contributor verification](../CONTRIBUTING.md#verification) for
deterministic seeds versus active fuzzing and the full qualification boundary.

Default limits are starting points, not universal safety policy. Lower them for
interactive or multi-tenant validation and measure real workloads before
raising them.

The versioned [threat model](threat-model.md) records current boundaries,
accepted risks, and the published v2 replacement-work repair. The original
v1 module through v1.1.1 does not independently admit replacement templates
or account for all consumed-name work; its subject and output limits do not
supply those missing bounds. See
[replacement admission advisory](https://github.com/faustbrian/go-ecma-regexp/security/advisories/GHSA-q6mc-8c55-cvwf).
Upgrade to the separate `/v2` module using the [migration guide](migration.md).
Until migration, keep replacement templates application-controlled or apply
explicit small admission limits and restrict substitution forms before
forwarding them to the legacy API.

Three exact `G115` source dispositions cover nonnegative, token-bounded
capture-count conversions and the pinned Unicode table sentinel encoding.
They do not exclude files or other scanner rules. The library maintainer
must review them when capture counting, token consumption, or generated
Unicode inputs change; current aliases select tables 0 through 434 and encode
them as values 1 through 435 in a `uint16` field.
