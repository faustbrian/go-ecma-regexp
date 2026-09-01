# ECMAScript regular-expression conformance matrix

The root module implements only the decision-bounded ECMA-262, Test262, Unicode, and JSON Schema profiles declared in the [specification decision register](../docs/specification-decisions.md). No broader JavaScript host-object or certification claim is implied.

[`manifest.json`](manifest.json) pins generated and official inputs. [`monitoring.json`](monitoring.json) records bounded source and change-authority review. [`decision-history.json`](decision-history.json) retains every canonical decision digest.

## Decision conformance

| Decision | Authority | Evidence boundary | Maintained-peer result |
| --- | --- | --- | --- |
| ECMAREGEXP-DEC-001 | `ecma262-source` | TestEditionIsExplicitAndClosed, TestTest262RegExpFeatureAccounting, TestParseRejectsMalformedAndUnsupportedSyntaxExplicitly | Not separately assessed |
| ECMAREGEXP-DEC-002 | `ecma262-source` | TestTest262RegExpFeatureAccounting, TestTest262RegExpSemantics, TestOverlappingLibraryDifferential | Registered overlapping vectors agree |
| ECMAREGEXP-DEC-003 | `ecma262-source` | TestFindReportsUTF16RuneAndByteIndices, TestUTF16InputViewMapsScalarAndUnpairedSurrogateBoundaries, TestDifferentialMatchingAgainstJavaScriptEngines | Registered overlapping vectors agree |
| ECMAREGEXP-DEC-004 | `ecma262-source` | TestInvalidUTF8MapsEachByteToReplacementCharacter, TestInputViewMapsUTF8BoundaryWidths | Not separately assessed |
| ECMAREGEXP-DEC-005 | `unicode-source` | TestUnicodePropertyEscapesUsePinnedUnicodeVersion, TestUnicodePropertyAliasesAreExact, TestUnicodeSetsPropertiesOfStrings, TestGenerateIsFormattedAndDeterministic | Not separately assessed |
| ECMAREGEXP-DEC-006 | `ecma262-source` | TestFlagsRejectDuplicatesConflictsAndUnknownFlags, TestAnnexBIsExplicitAndUnicodeModesRemainStrict, TestUnicodeSetsRejectMixedOperatorsAndReservedPunctuation | Not separately assessed |
| ECMAREGEXP-DEC-007 | `ecma262-source` | TestDuplicateNamedCapturesInDisjointAlternatives, TestDuplicateNamedCapturesThatMightBothParticipateAreRejected, TestDuplicateNamedCaptureReplacementUsesParticipatingGroup | Registered overlapping vectors agree |
| ECMAREGEXP-DEC-008 | `ecma262-source` | TestSessionImplementsGlobalLastIndex, TestSessionImplementsStickyLastIndex, TestFindAllAdvancesEmptyMatchesByUnicodeCodePoint, TestFindAllAndReplaceIncludeTheFinalUTF16Boundary | Not separately assessed |
| ECMAREGEXP-DEC-009 | `jsonschema-source` | TestJSONSchemaPatternIsUnicodeAndUnanchored, TestJSONSchemaPatternHonorsExplicitAnchors, TestJSONSchemaPatternRejectsNonUnicodeLegacySyntax | Not separately assessed |
| ECMAREGEXP-DEC-010 | `ecma262-source` | TestHostileExecutionPathsAreBounded, TestExecutionDoesNotLeakGoroutinesOrBuffers, TestMatchEnforcesWallTimeWithoutGoroutine, TestOperationsEnforceResultAndOutputLimits | Not separately assessed |
| ECMAREGEXP-DEC-011 | `test262-source` | TestTest262RegExpFeatureAccounting, TestTest262BuiltInNegativeRegExpLiteralSyntax, TestTest262RegExpLiteralNegativeSyntax, TestTest262RegExpSemantics | Not separately assessed |
| ECMAREGEXP-DEC-012 | `test262-source` | TestDifferentialMatchingAgainstJavaScriptEngines, TestOverlappingLibraryDifferential, TestUnicodeBackreferenceUsesSimpleCaseFolding | One classified JavaScriptCore policy divergence; all other registered vectors agree |
| ECMAREGEXP-DEC-013 | `ecma262-source` | TestReplaceImplementsECMAScriptSubstitutions, TestReplacementTokenBoundaries, TestSplitInsertsDefinedAndUnmatchedCaptures, TestUTF16OperationsPreserveExactInput | Not separately assessed |

## Official corpus and generated-data accounting

The pinned RegExp-literal slice accounts for all 186 negative files and all 52
positive files. Of the negative files, 167 exercise Pattern parsing and 19 are
JavaScript source-tokenization-only. Of the positive files, 12 execute matcher
calls and 40 exercise syntax or host evaluation only.

The pinned `built-ins/RegExp` inventory accounts for all 1,868 files: 192
negative Pattern cases, 443 generated class and property fixtures, 554 selected
matcher-feature files (494 executed and 60 without matcher calls), and 679
JavaScript host-object files outside this package's public surface. Exact rows
and executable bindings remain in [`conformance/test262.tsv`](conformance/test262.tsv).

Unicode code-point, case-folding, identifier, emoji-sequence, and
properties-of-strings tables are generated only from the digest-pinned Unicode
16.0.0 inputs in [`manifest.json`](manifest.json). The host Go or JavaScript
runtime Unicode version is not substituted.

[`conformance/differential.tsv`](conformance/differential.tsv) records exact
runtime or module versions, engine families, vector counts, and the classified
JavaScriptCore long-s divergence. Node.js and Deno share V8 and therefore count
as one engine family; Bun supplies JavaScriptCore, while regexp2 and goja cover
the smaller maintained-Go-peer overlap.

## Update and conformance process

1. A changed authority digest, release page, corpus pin, or maintained-peer result blocks unattended adoption.
2. Maintainers review the decision, requirement strength, compatibility and wire effects, executable evidence, and changelog digest before selecting behavior.
3. `golib specification check` validates the offline register, history, evidence, provenance, and change-control contract.
4. `golib specification check --online` verifies current authority content over bounded public HTTPS retrieval.
5. Test262 and maintained-peer gates remain separately attributable; peer agreement never overrides normative prose.

No unresolved decision is currently recorded.
