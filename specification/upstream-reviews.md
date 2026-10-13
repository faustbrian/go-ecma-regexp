# Upstream authority review history

This append-only record classifies discovery-authority changes before their
reviewed digests change. It does not change the selected specification edition.

## 2026-10-13: ECMA-262 publication discovery

- **Authority:** `ecma262-releases`
- **URL:** https://ecma-international.org/publications-and-standards/standards/ecma-262/
- **Previous SHA-256:**
  `0b0eb4cdcad935ddcd4dbfe389e44bd82206fbbc0bbf9a7f3c6b5cda15e5757e`
- **Reviewed SHA-256:**
  `80d750ca0510f00572de56ead77d97a450b3f63cdbf43f940e58f5f411a945df`
- **Retrieved and reviewed:** 2026-10-13
- **Applicability:** `ECMAREGEXP-DEC-001` directly and the ECMA-bound pattern,
  UTF-16, Unicode-mode, capture, progress, resource and replacement decisions
  (`002`, `003`, `004`, `006`, `007`, `008`, `010`, `013`).
- **Disposition:** Discovery acknowledged without adopting another edition;
  the closed ECMA-262 16th Edition (June 2025) contract remains selected.

Three consecutive TLS retrievals produced the reviewed digest, each with
92,327 bytes. The mutable publication page now identifies ECMAScript 2026,
17th Edition (June 2026), and its published page metadata reports modification
on 2026-06-30. It continues to link the archived 16th Edition. The previous
page body was not retained, so this review cannot reconstruct its byte-level
delta or claim the change was only presentation metadata.

Independent retrieval of the selected 16th Edition June 2025 PDF reproduced
its existing SHA-256:
`50e932767f16f3d4bf03e9d94157885125eeff1c165a8d8c05a867bcb56ee7da`
(6,659,783 bytes). The immutable source URL, edition, decision bindings,
Test262 revision, Unicode source and runtime semantics remain unchanged.
The current publication's normative HTML/printable PDF notice is not authority
to substitute a newer edition for this package's frozen source selection.

Adopting ECMAScript 2026 would require a separate compatibility-reviewed
semantic and conformance update; this discovery acknowledgement makes no
claim of 17th Edition support. The existing global monitoring review date is
preserved because this record concerns one discovery-authority delta.
