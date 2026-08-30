# Changelog

All notable changes to this project will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project uses semantic versioning.

## [Unreleased]

### Changed

- Replace copied repository tooling with the pinned go-library-tools v1.0.13
  contract while retaining package-owned policy and verification evidence.

### Documentation

- Declare the complete ECMA-262, Test262, Unicode, and JSON Schema
  interpretation surface in the
  [specification decision register](docs/specification-decisions.md), with
  strict machine records, conformance bindings, source monitoring,
  maintained-peer differential evidence, and append-only content history:
  - ECMAREGEXP-DEC-001 sha256:06fc6bec563a10c4bb86eaa1d9567f027b5ee86b074e05d69d492b0142557a4c
  - ECMAREGEXP-DEC-002 sha256:d7ba3054e0cf4a0c1ddefbe02c0edc6ed56e384724e25790acb8521c5751e954
  - ECMAREGEXP-DEC-003 sha256:15e79d2f5cac76a854210f68a2fd62cb09b788c91fef58e00035c272f8c2c7a6
  - ECMAREGEXP-DEC-004 sha256:129018c819bbd1cc16ce293427f198b1bca66126696509d9dad5be34e4f44ef0
  - ECMAREGEXP-DEC-005 sha256:1aeffca94e50698c7e2d7139ed7913565eb4643a172127a89582018c1eba4582
  - ECMAREGEXP-DEC-006 sha256:76a3c8601c5ae3a452faa67b4cd8d7dd5f35a4b85afa009f4ad20dd4714831e4
  - ECMAREGEXP-DEC-007 sha256:a8826a4999bf642e94eedad24d4e09c261507d421e6fa272316535045ea6a4e9
  - ECMAREGEXP-DEC-008 sha256:adc13a31805699778029f0fd0bf75c3b8469232f34616a518d9248928e902ecf
  - ECMAREGEXP-DEC-009 sha256:5d0a91fb12b8bee52c980a9af4fb34378546d795a9ea19b0ce0256d9bf89289c
  - ECMAREGEXP-DEC-010 sha256:904a373c1a9eaf0c9377af0d017158712b218c81ba81a56a0782d5f7a845424d
  - ECMAREGEXP-DEC-011 sha256:f32b68e9205973d69582df04974e363e8707b3221cea4fcabee34182a2615304
  - ECMAREGEXP-DEC-012 sha256:0333f7e28955d7a4adef530bdf23a56c416dc035fe462b3065e8fe9aa2ac194d
  - ECMAREGEXP-DEC-013 sha256:b1bd205e6b1f22094011061bf515005e428c145a87209eceebb8308080aafd37

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
