# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected module's docs.
2. Run standalone `go test ./...` or a focused affected-package test.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update the affected
stable entries in the [specification decision register](docs/specification-decisions.md),
and complete the Specification Decisions section of the pull request template.
An unresolved interpretation or stale source pin is release-blocking; peer
behavior cannot silently select policy.

Full scheduled and release mutation gates must finish with zero surviving viable
mutants. Local mutation is optional, not a universal contribution prerequisite.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for development and CI; the module minimum remains Go 1.27.0.
Standalone `go test ./...` runs deterministic tests, fuzz seeds and saved corpus
regressions without sibling checkouts or ecosystem ledgers. It does not actively
explore new fuzz inputs. External Test262 tests skip without `TEST262_ROOT`;
JavaScript differential tests use available engines and skip without enough
peers. Hosted replacement diagnostics are separate selected assertions, not
proof supplied by an ordinary local test run.

With the checksum-pinned public verifier selected by `.golib.yaml`, use:

```bash
make inventory
make check
make ci
```

Both Make routes use bounded local checks; the module declares no services or
test tags. Base-free local checks run deterministic fuzz seeds/corpora;
changed-source CI selects bounded active smoke for causally affected targets.
Ordinary CI consumes the public verifier binary. One existing race assertion
checks shared immutable programs through a synchronized caller-owned cache:
32 joined workers perform 100 finite matches each and assert the matched value.
It runs once with a two-minute timeout, not a repeated stress campaign.
Stateful sessions remain caller-owned and require external synchronization.

Short active exploration selects one target; optional local mutation selects
one package:

```bash
GOMAXPROCS=2 go test -run '^$' -fuzz '^FuzzCompileAndMatch$' -fuzztime=2s -parallel=1 -timeout=2m .
golib mutation --module . --package .
```

`make specification-check` validates the register offline;
`make specification-check-online` explicitly refreshes reviewed source digests.
`make check-full` retains online specification validation, hostile execution,
the ten-count leak assertion, all four native fuzz budgets, Test262/provenance,
documentation and benchmarks. Full interoperability retains Deno 2.9.3 and the
pinned V8/JavaScriptCore peer requirements, download checksums and cleanup.
Parser, matching, replacement, cancellation and limit changes retain their
affected fixtures and conformance; public API changes retain compatibility and
direct-consumer checks. Release qualification remains
`golib release dry-run --module .`. Select full campaigns for their risk or
release boundary. Report unavailable, skipped and failed checks honestly; local
tests are not external-corpus, hosted diagnostic or public-release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
