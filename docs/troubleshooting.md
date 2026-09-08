# Troubleshooting

## Compilation fails

Inspect the returned error rather than matching its text. Invalid flags and
syntax are distinct from [`LimitError`](api.md), which means the configured
parse or compile allowance was exhausted. Confirm that the pattern targets
ECMA-262 16th edition; later-edition syntax is intentionally rejected.

## A valid pattern does not find a match

`Program.Match` attempts only at `MatchOptions.StartUTF16`. Use `Program.Find`
for an unanchored search. A sticky (`y`) program still attempts only at the
explicit start position, and JSON Schema callers should use
`CompileJSONSchemaPattern` for Draft 2020-12 search semantics.

## Matching stops with a limit, timeout, or cancellation error

The default budgets are finite. Treat `LimitError`, `TimeoutError`, and context
cancellation as execution failures, not ordinary non-matches. Adjust limits
only after measuring bounded production inputs; see [limits and
security](security.md) and [performance](performance.md).

## Reported indices differ from Go byte offsets

ECMAScript indices use UTF-16 code units. Read `Index.UTF16` for the normative
position and use `Index.Byte`, `Index.Rune`, and `Index.Exact` when mapping to a
Go string. Use the `UTF16String` APIs when lone surrogates must be preserved.

## More help

Use [GitHub Issues](https://github.com/faustbrian/go-ecma-regexp/issues) for a
reproducible defect and [GitHub
Discussions](https://github.com/faustbrian/go-ecma-regexp/discussions) for
adoption questions. Report suspected vulnerabilities only through the private
process in [SECURITY.md](../SECURITY.md).
