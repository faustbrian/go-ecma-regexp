# API and indices

The processing pipeline is explicit:

1. `Tokenize` returns immutable tokens with pattern byte spans.
2. `Parse` returns an immutable typed `Pattern` AST.
3. `Compile` returns an immutable executable `Program`.
4. `Match`, `Find`, `FindAll`, `Replace`, and `Split` execute with caller
   options and context cancellation. A canceled context is checked before
   input admission; preparation and VM execution poll cancellation
   cooperatively. Nil contexts retain background behavior.

`Match` attempts only `MatchOptions.StartUTF16`. `Find` searches from that
position unless the `y` flag makes the program sticky. Starts must be nonnegative
UTF-16 positions: existing Unicode Find-based paths can panic on negative values.
`Split` starts at zero and ignores `StartUTF16` and `g`/`y` selection.
`FindAll` returns
ordered, non-overlapping matches and advances empty matches using ECMAScript
`AdvanceStringIndex` semantics.

`Replace` and `ReplaceUTF16` admit the subject and replacement independently
against `MatchLimits.InputBytes` and `InputRunes`, even without a match. Exact
UTF-16 storage costs two bytes per unit; paired surrogates count as one code
point and lone surrogates as one. Token scans and capture-name conversion share
`Steps` with the VM; output allowances do not replace input or work admission.
This is the published v2 admission contract; the former v1 subject-only policy
requires explicit major-module migration.

Every capture has a half-open `IndexSpan`. `Index.UTF16` is the normative
ECMAScript code-unit position. `Index.Rune` and `Index.Byte` map to the Go
input. `Index.Exact` is false at a boundary inside a surrogate pair. For exact UTF-16
input containing any lone surrogate, every post-origin Go rune/byte mapping is
unavailable (`-1`, `Exact=false`), including otherwise scalar boundaries. Check
participation before using a capture span; unmatched captures have zero metadata.
Only a successful matched result has a usable `Full` capture; the zero Result panics.

The string-taking APIs decode invalid UTF-8 one byte at a time as U+FFFD and
preserve byte offsets. Use `UTF16String` APIs when exact ECMAScript strings,
including lone surrogates, are required. `GoString` rejects an unpaired
surrogate; `LossyString` performs explicit replacement.

`Program` values are immutable and may be shared concurrently. `Session`
provides caller-owned `lastIndex` state for `g` and `y`; it is not safe for
concurrent mutation. On an ordinary failed global or sticky execution,
`lastIndex` resets to zero. A limit or cancellation error leaves it unchanged.
Empty successful Session execution does not advance past that boundary: callers
looping Exec own progress. Collection APIs advance empty matches themselves.

AST and result slice/map/unit accessors return copies or immutable values. Original
pattern source remains retained by Pattern/Program and can be sensitive; captures
copy input units. No Unicode normalization is performed. `Node.Ranges` reports
stored non-builtin terms, not expanded property membership or flattened v-mode
algebra; property terms can have zero-valued endpoints.

Zero option fields do not acquire defaults. Use the explicit default constructors.
Parse/Compile have no context API. Steps charge dispatch and substitution/name
work rather than every internal comparison. Allocations count selected logical
reservations, not total Go allocations/bytes, preparation or successful capture
copies. WallTime starts after subject/template preparation and is sampled; caller
deadlines remain cooperative and do not interrupt every instruction or copy.
