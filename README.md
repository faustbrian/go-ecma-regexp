# ecma-regexp

[![CI](https://github.com/faustbrian/go-ecma-regexp/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-ecma-regexp/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-ecma-regexp/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-ecma-regexp.svg)](https://pkg.go.dev/github.com/faustbrian/go-ecma-regexp)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-ecma-regexp?sort=semver)](https://github.com/faustbrian/go-ecma-regexp/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`ecma-regexp` is a bounded, specification-backed ECMAScript regular
expression engine for Go. It implements JavaScript semantics directly; it
does not translate patterns to Go's RE2-based `regexp` package and does not
embed a JavaScript runtime.

The supported language is closed to ECMA-262, 16th edition (ECMAScript 2025),
with Unicode 16.0.0 data. The exact ECMA-262, Test262, Unicode, and emoji
provenance is recorded in [`specification/manifest.json`](specification/manifest.json).

> Release status: stable and released as v1.0.0. The applicable pinned Test262
> inventory, meaningful 100% production statement coverage, and scoped mutation
> gate are complete. See the
> [conformance inventory](specification/README.md).

## Quick start

```go
program, err := ecmascript.Compile(
	`(?<word>\p{Letter}+)`,
	"u",
	ecmascript.DefaultCompileOptions(),
)
if err != nil {
	return err
}

result, matched, err := program.Find(
	ctx,
	"42 Helsinki",
	ecmascript.DefaultMatchOptions(),
)
if err != nil {
	return err
}
if matched {
	word, _ := result.Named("word")
	fmt.Println(word.Value().LossyString())
}
```

All parsing and execution limits are explicit. Limit exhaustion, cancellation,
and wall-time exhaustion are errors distinct from invalid syntax and an
ordinary no-match result.

For JSON Schema Draft 2020-12, use `CompileJSONSchemaPattern`; it selects
Unicode semantics and the required unanchored search behavior.

## Documentation

- [Support and compatibility](docs/support.md)
- [Specification decisions](docs/specification-decisions.md)
- [Syntax, flags, Unicode, and captures](docs/syntax.md)
- [API and index semantics](docs/api.md)
- [Replacement behavior](docs/replacement.md)
- [JSON Schema profile](docs/json-schema.md)
- [Limits and security](docs/security.md)
- [Performance and benchmarks](docs/performance.md)
- [Migration from Go regexp and PCRE](docs/migration.md)
- [Cookbook](docs/cookbook.md)
- [FAQ](docs/faq.md)
- [Changelog](CHANGELOG.md)

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/README.md)
and its [Domain utilities family](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/design-language.md#package-families-and-selection).

## Development gates

```sh
make inventory
make cohesion
make check
make ci
```

The gates use the pinned `go-library-tools` release declared in `.golib.yaml`.
Package-specific safety, documentation, Test262 conformance, interoperability,
fuzz, benchmark, and mutation evidence remain declared or retained in this
repository.

The package uses no `unsafe`, hidden workers, or global mutable caches. A
compiled `Program` is immutable and concurrency-safe. Stateful `Session`
values and any application cache are caller-owned and require external
synchronization when shared.

## Documentation

Start with the [documentation index](docs/README.md) for syntax, replacement,
security, compatibility, and migration guidance.
