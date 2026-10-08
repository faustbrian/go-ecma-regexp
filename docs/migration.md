# Migration from Go regexp and PCRE

## Published v2 adoption

V2.0.0 is published from main. Use
`github.com/faustbrian/go-ecma-regexp/v2` imports and a v2 Git tag; source
remains in the repository root. Go 1.27 is required. Published v1.1.1 remains
available at its original module path and tag, without a backport of the
replacement admission repair. See [security guidance](security.md).

Subject and replacement templates independently use `InputBytes` and
`InputRunes`. UTF-16 templates cost two bytes per unit and one code point per
surrogate pair or lone unit. Budget both inputs explicitly; zero remains a
zero allowance. Consumed token/name scanning and conversion share `Steps`
with matching, even when the template produces no output. Previously accepted
templates can now return a typed limit error without partial output.

Caller cancellation is checked before admission and cooperatively during
preparation. Continue supplying a caller deadline for the whole operation;
VM `WallTime` does not include preparation. Session state remains unchanged
when execution is refused.

## Other engine migrations

Do not assume that a pattern accepted by Go `regexp`, PCRE, or ECMAScript has
the same language or result in another engine.

From Go `regexp`:

- replace `regexp.Compile` with `Compile` and provide an explicit flag string;
- choose `Match` for exact-start matching or `Find` for search;
- supply finite compile and match options plus a context;
- consume UTF-16 indices when reproducing JavaScript-visible behavior;
- audit backreferences, lookaround, named captures, Unicode properties, and
  replacement tokens instead of rewriting them to RE2 approximations.

From PCRE:

- remove PCRE-only options, verbs, conditionals, recursion, and syntax;
- verify escape and character-class behavior in the selected ECMAScript mode;
- audit anchors, newline handling, duplicate names, and case folding;
- use the differential and conformance tests as executable migration vectors.

For JSON Schema, migrate directly to `CompileJSONSchemaPattern`; do not add
implicit anchors or pass user-selected flags.
