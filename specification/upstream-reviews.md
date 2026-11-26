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

## 2026-11-26: release-index monitoring review

The two discovery authorities whose digests no longer matched their pins were
retrieved twice over HTTPS with the same request profile as the online
specification check. Each pair of retrievals produced the same digest. All
other monitoring entries, including the selected 16th Edition PDF and the
Unicode 16.0.0 UCD archive, still match their pins. The thirty-day cadence and
all authority URLs remain unchanged.

| Authority | Previous SHA-256 | Current SHA-256 |
| --- | --- | --- |
| `ecma262-releases` | `80d750ca0510f00572de56ead77d97a450b3f63cdbf43f940e58f5f411a945df` | `7a697cd489c446c25971c4f8e54209eb5dd514e247e9502f35084f9ce7e05ce3` |
| `unicode-releases` | `4c3532ee2bd1fba966a2f303e9d67b686c7ad782a08b9188ae722ec4904fd68f` | `8c59d14866ec43ff6973e0cede0954b86306466907a13fccfabb01a58fe419f4` |

The ECMA-262 publication page (92,775 bytes) still identifies ECMAScript 2026,
17th Edition (June 2026), as current, reports page modification on 2026-06-30
and continues to link the archived 16th Edition (June 2025). No new edition or
corrigendum is listed. The 17th Edition was already acknowledged on
2026-10-13 without adoption.

The page bytes change with the site's software. A public archive copy from
2026-11-18 differs from the current page
only in WordPress, theme, plugin and asset version strings and in the markup of
the search form. The copy behind the previous pin was not retained or found in
the archive, so the earlier delta is not attributed byte for byte.

The Unicode versions index (34,981 bytes) still identifies Unicode 18.0.0 as
the latest version, as acknowledged in the 2026-10-01 monitoring review. A
public archive copy from 2026-10-05 reproduces the previous pinned digest
exactly. Compared with it, the current page uses LF instead of CRLF line
endings, and its section on linking to the core specification now says that a
"latest" link drops the chapter part of the path, with the example changed
from `https://www.unicode.org/versions/latest/core-spec/chapter-2/#G4315` to
`https://www.unicode.org/versions/latest/core-spec/#G4315`. No version, date,
citation or data-file statement changed.

Neither change selects a new edition or data version: `ECMAREGEXP-DEC-001`
keeps the closed ECMA-262 16th Edition contract and `ECMAREGEXP-DEC-005` keeps
the generated Unicode 16.0.0 tables until an explicit, compatibility-reviewed
upgrade. Normative pins, generated
inputs, Test262 bindings, decisions and runtime semantics are unchanged.
