# Threat model

- Model: `ECMA-REGEXP-THREAT-MODEL-2.0.0`
- Applies to: `github.com/faustbrian/go-ecma-regexp/v2` candidate on main
- Reviewed: 2026-10-05
- Owner: ECMA regexp maintainers

## Assets and trust boundaries

The package protects process availability, matching and index correctness,
caller-owned execution state, and transient input lifetime. Patterns, flags,
input strings, replacement templates, start positions, and caller-selected
limits can originate outside the application trust boundary. Contexts and
configured budgets express the application's availability policy. Results,
captures, indices, and errors return to application code.

Compiled programs are immutable and may be shared. A Session is mutable and
requires caller-owned synchronization. Program values retain their pattern
source; applications must not compile secrets into patterns or log sensitive
patterns, inputs, captures, or replacement results.

Construction and execution introduce no ambient network, filesystem,
credential, process, worker-goroutine, or global-cache boundary. Repository
generators and interoperability tools are separate build-time boundaries with
pinned dependencies and specification inputs.

## Controls and cancellation

Pattern bytes, AST nodes and depth, captures, character classes, and program
instructions have finite parse and compile budgets. Execution separately
bounds input bytes and runes, VM steps, backtracking, stack and recursion,
logical allocations, results, output units, and VM wall time. Zero allowances
remain zero, not unlimited.

Replacement templates are admitted independently against the same input byte
and code-point allowances as subjects. UTF-16 storage costs two bytes per
unit; a surrogate pair counts as one code point and a lone surrogate as one.
Admission allocates no template copy or index map. Consumed token/name scanning
and name conversion share the VM step budget, including names producing empty
output and suffixes revisited during literal fallback. Temporary name storage
is bounded by admitted template bytes. Caller cancellation is checked during
admission and by periodic charged execution work.

Every public execution method checks genuine caller cancellation before input
admission. UTF-8 counting/construction and UTF-16 counting/copy/construction
poll cancellation in bounded chunks. Observed cancellation returns no partial
output; a Session retains its prior lastIndex when execution is refused.
Nil context retains background behavior. Context errors, LimitError, and
TimeoutError remain distinct from an ordinary non-match.

String inputs decode invalid UTF-8 one byte at a time as replacement runes;
exact UTF-16 inputs retain unpaired surrogates and explicit index mappings.
These semantic contracts and finite input limits are unchanged by the
cancellation repair.

## Verification boundaries

Hosted public regression tests independently exercise pre-cancellation and
expired-deadline precedence on all thirteen execution entry points, absent
partial output, and the out-of-range global Session state boundary. Ordinary
limit and nil-context controls remain separate. Corrected-source hosted
verification is required before claiming this repair delivered.

These tests do not establish runtime cancellation during a specific preparation
pass. The periodic checks are source-reviewed; no manufactured context polling
counts or timing-stress result is treated as that runtime proof. Existing
semantic, input-boundary, and index tests remain required affected checks.

## Accepted risks

| ID | Severity | Owner | Rationale | Mitigation | Review condition |
| --- | --- | --- | --- | --- | --- |
| ECMA-RISK-001 | Medium | Application integrator | Cancellation is cooperative; a bounded chunk can finish before the next check. | Finite input/VM budgets and periodic context checks; caller deadlines include preparation. | Revisit if a bounded chunk exceeds the application's latency budget or chunk size changes. |
| ECMA-RISK-002 | Medium | Runtime maintainers | Tokenization, parsing, and compilation have no context-taking API. | Finite structural and instruction budgets; applications choose stricter limits for interactive workloads. | Revisit for a context-aware compilation API or an observed unacceptable bounded compile delay. |
| ECMA-RISK-003 | Medium | Application integrator | VM wall time excludes input preparation and is not a whole-operation deadline. | Supply a bounded caller context and finite input budgets. | Revisit if wall-time scope changes or a whole-operation timeout is promised. |
| ECMA-RISK-004 | Medium | Application integrator | Caller-selected large budgets increase resource exposure. | Tenant-specific finite limits, bounded program caches, and isolated process capacity where needed. | Revisit when raising defaults or supported limits, or after a resource-exhaustion incident. |

## Replacement repair and release boundary

ECMA-OPEN-001 is owned by runtime maintainers. Hosted pre-fix regression tests
demonstrate missing replacement admission and uncharged consumed-name work.
The candidate repair admits templates and charges scanning/conversion, but
corrected-source hosted verification and independent review remain pending.
This finding is not yet closed or accepted. The narrowed replacement contract
requires a major release and documented migration from subject-only budgets;
v1 patch publication does not qualify the repair. Runtime interruption during
a particular scan remains unproved; source checkpoints are not that proof.

Review this model after changes to input handling, budgets, cancellation,
concurrency, dependencies, generated data, or public execution semantics.
