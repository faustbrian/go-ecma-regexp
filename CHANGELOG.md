# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project uses semantic versioning.

## [Unreleased]

### Changed

- Replace copied repository tooling with the pinned go-library-tools v1.0.4
  contract while retaining package-owned policy and verification evidence.

### Documentation

- Replace the archived monorepo link with package-owned documentation.

## [1.0.0] - 2026-08-25

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Link the package README to package-owned documentation.

### Fixed

- Express UTF-16 narrowing bounds directly so static analysis can verify that
  parsed code points and literal units cannot truncate.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-ecma-regexp` identity while preserving its documented API and behavior.
- Provision Test262 in a task-owned disposable directory for conformance and
  provenance gates instead of relying on shared `/tmp` state.
- Delegate local mutation checks to the canonical exact-100 repository runner
  instead of maintaining package-local exclusions and reduced thresholds.
- Mutation campaigns reuse the external Test262 coverage baseline while
  exercising each mutant without reinvoking the external corpus.
- Provision pinned Node and Deno runtimes for the required multi-engine
  differential test in repository CI.

### Distribution

- Include the canonical MIT licence in the independently published module.

### Compatibility

- Added a pinned module export baseline so incompatible public API changes
  fail the canonical repository gate.

### Added

- A canonical, evidence-linked specification decision register covering the
  supported edition, Pattern-only scope, text and index models, Unicode modes,
  JSON Schema behavior, bounded execution, Test262, and peer interoperability.
- Closed ECMAScript 2025 tokenizer, parser, immutable AST, compiler, and
  bounded matcher.
- Captures, backreferences, lookaround, scoped modifiers, Annex B grammar,
  Unicode 16.0.0 properties, and Unicode Sets.
- UTF-16-exact matching, search, replacement, split, and stateful execution.
- JSON Schema Draft 2020-12 pattern profile.
- Differential Node.js and Deno vectors plus fuzz surfaces.
- Complete applicable pinned Test262 accounting with delegated matcher runs
  and structural generated-data verification.
- A canonical interoperability target that provisions and verifies the pinned
  Test262 corpus before running conformance and differential checks.
- Separate official Test262 conformance from independent-engine differential
  interoperability so both results remain attributable.

### Security

- Parse decimal backreferences and hexadecimal escapes at their destination
  widths before conversion, and branch before converting UTF-16 code units.
- Finite parse, compile, execution, output, and wall-time budgets.
- Synchronous context cancellation without hidden worker goroutines.
- Meaningful 100% production statement coverage and exact 100% mutation
  efficacy for every viable mutant.

[Unreleased]: https://github.com/faustbrian/go-ecma-regexp/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/faustbrian/go-ecma-regexp/releases/tag/v1.0.0
