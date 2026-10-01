# Support and compatibility

## Version matrix

| Surface | Supported value |
|---|---|
| ECMA-262 | 16th edition, ECMAScript 2025 |
| Unicode | 16.0.0 |
| JSON Schema | Draft 2020-12 regular-expression profile |
| Go | version declared by `go.mod` |
| Core runtime | Pure Go; no JavaScript runtime |

Future ECMAScript syntax is rejected until its grammar and semantic changes
are inventoried, implemented, and evidenced. Earlier editions are not exposed
as approximate aliases.

## Feature matrix

| Feature | Status |
|---|---|
| Literals, alternation, groups, quantifiers | Implemented |
| Assertions and word boundaries | Implemented |
| Numbered and named captures | Implemented |
| Backreferences | Implemented |
| Lookahead and lookbehind | Implemented |
| `d g i m s u v y` flags | Implemented |
| Unicode properties and case folding | Unicode 16.0.0 tables |
| Unicode Sets intersection, subtraction, strings | Implemented |
| Annex B pattern grammar | Configurable; enabled by default |
| Replacement and split semantics | Implemented |
| `lastIndex` execution | Caller-owned `Session` |
| Canonical pattern formatter | Not exposed |
| Global compiled-pattern cache | Intentionally not provided |

The implementation status is not a conformance certification. The exact
applicable Test262 execution and skip boundary is recorded in the specification
inventory; JavaScript RegExp object APIs are outside this package's surface.
The differential gate runs against pinned Node.js, Deno, and Bun releases so
that both V8 and JavaScriptCore behavior is visible. Engine-specific exceptions
are enumerated in `specification/conformance/differential.tsv`.

Optional integrations should depend only on the public `Compile`, `Program`,
and JSON Schema profile APIs. The core module does not require
`json-schema`, `rule-engine`, or `validation`.

The owned interoperability operation provisions checksum-pinned Node 24.4.1
(V8) and Bun 1.3.14 (JavaScriptCore) in disposable storage and requires both
families before comparing the existing vectors and maintained Go peers.
It also compares Deno 2.9.3 when present. Ordinary suite runs retain their
available-runtime comparisons; they alone are not release-family evidence.
Run `make -f verification/package.mk interoperability` from the repository
root for this same gate. The peers are test-only and never runtime dependencies.

The October 1, 2026 monitoring review found the ECMA release page still
publishes the 17th edition (ECMAScript 2026), with the 16th-edition archive
retained. The Unicode release index now identifies Unicode 18.0.0 as latest.
Neither notice selects a new package edition or Unicode data version:
ECMAScript 2025 and Unicode 16.0.0 remain the supported semantic boundary.
The prior index bodies were not retained, so the Unicode change is not
classified as presentation-only.
