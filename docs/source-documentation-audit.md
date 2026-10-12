# Source documentation audit

## Pre-edit scope and decisions

Base `d070f369126227dcb8470962beaf2de28e153d4b`; docs/comments-only Tier A.
One active root module `github.com/faustbrian/go-ecma-regexp/v2`, three source
directories, 63 tracked Go files. Root is the only public package; two internal
commands own Test262 bridge and Unicode generation, not application runtime.

Reproduce with `git ls-files '*.go'`, `git ls-files '*go.mod'`, the module/package
catalogs, and `golib source-docs inventory --json` / `source-docs check --json`.
Initial report: 157 missing useful public/member/command comments plus two
local-generator availability findings on recognized generated output. These
counts locate gaps, not acceptance or specification conformance proof.

## Complete tracked Go inventory

- `.`: `a_matcher_anchor_internal_test.go`, `annexb_test.go`, `ast.go`, `benchmark_test.go`, `branches_internal_test.go`, `capture_semantics_test.go`, `casefold_test.go`, `compiler.go`, `compiler_boundary_internal_test.go`, `contract_internal_test.go`, `deterministic_execution_boundaries_diagnostic_test.go`, `differential_gate_test.go`, `differential_test.go`, `duplicate_names_test.go`, `edition.go`, `edition_test.go`, `errors.go`, `example_test.go`, `execution_cancellation_diagnostic_test.go`, `flags.go`, `fuzz_test.go`, `grammar_match_test.go`, `hardening_test.go`, `input.go`, `input_checkpoint_test.go`, `input_contract_internal_test.go`, `input_copy_diagnostic_test.go`, `jsonschema.go`, `jsonschema_test.go`, `library_differential_test.go`, `limits_internal_test.go`, `lookaround_test.go`, `matcher.go`, `matcher_boundary_internal_test.go`, `matcher_test.go`, `modifier_test.go`, `mutation_resource_guard_unix_test.go`, `mutation_small_contract_test.go`, `named_identifier_test.go`, `operations.go`, `operations_boundary_internal_test.go`, `operations_test.go`, `parser.go`, `parser_boundary_internal_test.go`, `replacement_admission_diagnostic_test.go`, `replacement_unicode_diagnostic_test.go`, `session.go`, `session_test.go`, `span.go`, `syntax.go`, `syntax_test.go`, `test262_generated_test.go`, `test262_semantics_test.go`, `test262_syntax_test.go`, `unicode.go`, `unicode_property_test.go`, `unicode_sets_test.go`, `unicode_tables_generated.go`, `utf16.go`, `utf16_test.go`.
- `internal/cmd/test262bridge`: `main.go`.
- `internal/cmd/unicodegen`: `main.go`, `main_test.go`.

## Contract and private rationale coverage

- Edition2025/ECMA-262 16th edition, validated flag alphabet and Annex B legacy
  profile, Unicode 16.0.0 exact aliases/folding and explicit JSON Schema 2020-12
  Unicode/unanchored profile; no JavaScript host objects, coercion, callbacks,
  PCRE fallback or Go regexp delegation. Interpretations/authority pins and
  dated specification decisions remain frozen, not replaced by peer behavior.
- Token source byte spans; parsed immutable nodes/capture metadata, duplicate
  name disjointness and backreference resolution; Unicode Sets algebra/string
  properties and legacy-vs-u/v escape rules. Compiler reversal for lookbehind,
  capture clearing per quantified iteration, longest class-string preference,
  snapshot ownership and progress guards are the non-obvious private seams.
- Synchronous backtracking VM: preferred branches, copied capture/guard frames,
  positive/negative lookaround capture ownership, shared attempt counters,
  sampled cancellation/wall time and charged work. No linear-time, process-wide
  memory, total instruction-cost or hard-interruption guarantee is inferred.
- Exact UTF16String copies, invalid UTF-8 replacement, code-unit versus code-point
  matching, capture participation, normative UTF-16 and optional Go coordinates;
  input preparation admission/polling and lone-surrogate mapping limitations.
- Match/Find and collection/replace/split entrypoints; explicit starts, sticky/global
  variants, shared operation budgets, zero-width progress and no partial outputs
  on reported errors. Session owns mutable lastIndex and external synchronization;
  errors preserve it, ordinary stateful non-match resets it, empty success does
  not automatically advance state. Programs/results/AST values expose copies.
- All test sources inventoried: parser/compiler/VM limit and branch matrices,
  Annex B/modifiers/casefold/Unicode/captures/session/operation boundaries, input
  and replacement hosted diagnostics, Test262 semantics/syntax/features, peer and
  overlapping-library differential variants. Four fuzz targets, four benchmarks
  and one public Example. Platform-specific mutation resource guard owns its lock,
  stop channel and diagnostic workers; it is not production engine behavior.
- Unicode generator validates three pinned source digests and emits formatted
  deterministic tables; bridge handles one synthetic JSON/stdin request, runs
  finite Session execution, reports one stdout response and is not a public API.
  The generated Unicode table (75,848 lines) and all oracle/fixture authorities
  remain byte-exact. test262_generated_test.go has no generated header: it is a
  maintained external-corpus test variant, not newly regenerated source.

## Missing comments versus existing limitations/questions

Add useful field/method/enum and command documentation and private rationale;
keep expressions, API baseline, dependencies, compiler/analyzer directives and
Unicode producer/output unchanged. The generator header names the existing local
`internal/cmd/unicodegen` directory; `unicode.go` has the existing go:generate
script and scripts/sync-unicode.sh invokes that producer with digest-pinned inputs.
The initial shared checker reported generator-review and generator-unavailable
at output line3 despite Generated=true. Coordinator repaired local generator-directory
recognition in the shared tool; this repository output/producer/directive did not change.
The remaining shell regeneration declaration is reviewed manually for availability
and pinned output ownership, without executing downloads or arbitrary shell.

Truthful limits: Steps charge dispatch/substitution/name work, not each comparison
inside backreferences/classes. Allocations charge selected logical reservations,
not all heap allocations or capture/input/output copies. WallTime starts after
input and replacement preparation and is sampled; completed copies and individual
instruction work are not asynchronously preempted. Parse/Compile have no context.
AST nodes and instruction caps do not directly cap all temporary compile storage
or repetition expansion work. Empty repeated bodies may emit no instructions;
overlarge numeric counts saturate to platform max int. This is an existing
compile-control limitation observed in source, not runtime-probed or repaired.

StartUTF16 must be nonnegative. Existing Unicode Find-based paths index the
boundary table before validating negatives and can panic; no runtime repair is
bundled. Match attempts one index; Split ignores StartUTF16 and global/sticky
selection. Any lone surrogate in exact input makes every post-origin Go byte/rune
mapping unavailable (-1/Exact=false), not merely that surrogate boundary. Unmatched
captures have zero values/spans; only successful results have a usable Full capture.
Node.Ranges is stored metadata, not expanded Unicode property/class-set membership.

Unknown rationale/behavior decisions: negative-start handling consistency across
entrypoints, fuller Go boundary mappings for exact strings containing lone
surrogates, and whether individual instruction work should have finer charging.
These are separate behavior/security qualification questions, not invented fixes
or new claims of official conformance.

## Selected checks

Request one bounded serialized Go1.27.2 session with owned GOCACHE/GOMODCACHE/GOTMPDIR:
format only changed Go; exact comment-free AST/signature/raw token/literal identity
for all63 tracked files including generated data, plus exact directives; parsed
comments, root and two command GoDoc, native scripts/check-docs.sh (Example plus
vet). `.golib.yaml` uses apidiff api/baseline.txt, not a full-prose snapshot;
retain baseline bytes. Do not rerun Test262, engine peers, race/leak, active fuzz,
mutation or performance campaigns for unchanged algorithms.

Also check affected maintained prose/local links, source inventory/check, generator
availability structurally without executing downloads, offline specification
register integrity, markers, pinned added/staged-content secrets, Git whitespace,
one ordinary complete final-diff Tier A review and explicit Intent/ACK commit.
Remote/main/CI, publication, enforcement tooling and aggregate ledger belong to
coordinator. Cache/helper/probe resources are cleaned on success and failure.

## Results

- Go1.27.2 bounded native proof: gofmt clean across all63 tracked Go files;
  exact raw tokens/literals/implicit semicolons, comment-free AST/signatures and
  compiler/analyzer directives equal to recorded base. No token normalization
  exception was needed. All245 comment groups parsed/rendered, including generated
  table metadata. Only changed hand-maintained source was formatted.
- Root and both internal-command native `go doc -all` outputs rendered and were
  inspected. Apidiff mode applies, so no full-prose snapshot is regenerated.
- `make -f verification/package.mk docs` passed its structural/link checks,
  existing public Example (`go test . -run '^Example' -count=1`) and `go vet ./...`.
- Two tiny disposable observations passed: Unicode Find with a negative start
  panics, and a scalar capture after a lone surrogate has unavailable Go mappings.
  These confirm existing limitations, not repairs or specification interpretations.
  Potential excessive empty-repetition expansion is source-observed only; no
  hostile expansion/runtime conformance or performance probe was executed.
- Updated shared source checker: exit0/errors0 with one legitimate generator-review
  at unicode_tables_generated.go:3. Manual producer/script/pin availability review
  completed without expanding arbitrary shell or regenerating Unicode data.
  Offline specification check passed all13 decisions; frozen authorities unchanged.
- Exact allthree owned caches, helper/render/inventory evidence and in-worktree
  probe are removed after evidence capture, including downloaded readonly module
  directory permission normalization only inside the owned cache. Shared tooling
  and the detached worktree remain preserved.

Final affected-link/marker, pinned staged-added secret, whitespace and ordinary
complete TierA final-diff results are reported with the Intent commit handoff.
Unverified: full Test262 and peer interoperability, timing/concurrency/leak/fuzz/
mutation/performance campaigns, physical heap ceilings or hard interruption,
external consumers, hosted CI and remote-main delivery. This prose pass does not
repair the existing runtime limitations or qualify those boundaries.

## Objective Enforcement

The checksum-pinned public Tools v2.4.1 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.

This patch admits the existing tracked Go command-package producer directory.
The generated Unicode table, header, shell directive, producer and pinned
inputs are unchanged; unsupported shell regeneration still requires manual
review, not an assertion that regeneration was executed.
