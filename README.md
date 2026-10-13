# ecma-regexp

[![CI](https://github.com/faustbrian/go-ecma-regexp/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-ecma-regexp/actions/workflows/ci.yml)
[![Documentation](https://img.shields.io/badge/docs-adoption_and_API-blue)](docs/index.md)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-ecma-regexp/v2.svg)](https://pkg.go.dev/github.com/faustbrian/go-ecma-regexp/v2)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-ecma-regexp?sort=semver)](https://github.com/faustbrian/go-ecma-regexp/releases)
[![Go](https://img.shields.io/badge/go-1.27.0-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`ecma-regexp` is a bounded, specification-backed ECMAScript regular
expression engine for Go. It implements JavaScript semantics directly; it
does not translate patterns to Go's RE2-based `regexp` package and does not
embed a JavaScript runtime.

The supported language is closed to ECMA-262, 16th edition (ECMAScript 2025),
with Unicode 16.0.0 data. The exact ECMA-262, Test262, Unicode, and emoji
provenance is recorded in [`specification/manifest.json`](specification/manifest.json).

> Release status: the published v2 major is stable and requires Go 1.27
> or later. The legacy v1 module has no backport of the replacement
> admission repair. See [security guidance](docs/security.md),
> [migration](docs/migration.md), and the
> [conformance inventory](specification/README.md).

## Adoption entry points

The root module is `github.com/faustbrian/go-ecma-regexp/v2` (stable).
It requires Go 1.27.0. The prepared `v2.0.1` patch restores negative-start
admission and bounded empty-repetition compilation. After its tag is published,
adopt it with `go get github.com/faustbrian/go-ecma-regexp/v2@v2.0.1`.
Root subpackages share this module; they need no separate installation.

Start with the [documentation map](docs/index.md) for package selection,
examples, operational constraints, and contributor references.

## Install

```sh
go get github.com/faustbrian/go-ecma-regexp/v2@v2
```

Use `/v2` imports for the current stable major. Source remains in the same
root directory on main, without a version-specific source branch or
directory. See [migration guidance](docs/migration.md).

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
ordinary no-match result. Budgets count selected logical work, not every internal
comparison or all heap allocations; execution time/cancellation is cooperative.
Use nonnegative UTF-16 start indices and application-owned capacity limits.
See [the exact API boundaries](docs/api.md).

For JSON Schema Draft 2020-12, use `CompileJSONSchemaPattern`; it selects
Unicode semantics and the required unanchored search behavior.

The checked-in [`Example`](example_test.go) is the compiler-checked public
quick start and runs as part of the documentation gate.

## Package map

| Import path | Package | Use |
| --- | --- | --- |
| `github.com/faustbrian/go-ecma-regexp/v2` | `ecmascript` | Parse, compile, and execute bounded ECMAScript and JSON Schema regular expressions. |

The module exposes no public subpackages. Commands beneath `internal/cmd` are
repository tooling and cannot be imported by consumers.

## Documentation

Start with the [documentation index](docs/README.md) for the complete guide
set.

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
- [Troubleshooting](docs/troubleshooting.md)
- [Support](SUPPORT.md)
- [Security reporting](SECURITY.md)
- [Changelog](CHANGELOG.md)

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/README.md)
and its [Domain utilities family](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/design-language.md#package-families-and-selection).

## Development gates

```sh
go test ./...
```

`make check` and `make ci` use bounded local checks with the pinned public
verifier declared in `.golib.yaml`; `make check-full` retains explicit full
assurance. Ordinary tests run fuzz seeds, not active exploration, and do not
prove external Test262 or release interoperability when their peers are absent.
See [contributor verification](CONTRIBUTING.md#verification) for focused fuzzing,
optional local mutation and selected conformance/release requirements.

The package uses no `unsafe`, hidden workers, or global mutable caches. A
compiled `Program` is immutable and concurrency-safe. Stateful `Session`
values and any application cache are caller-owned and require external
synchronization when shared.
