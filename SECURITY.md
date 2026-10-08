# Security Policy

## Reporting

Report suspected vulnerabilities privately through the GitHub security
advisory for `faustbrian/go-ecma-regexp`. Do not open a public issue containing exploit
details, credentials, private fixtures, or affected deployment information.

Include the affected module and version, impact, reproduction, preconditions,
and any suggested mitigation. Reports are acknowledged as soon as practical;
timelines depend on severity and verification.

## Supported Versions

The latest stable major is
[`github.com/faustbrian/go-ecma-regexp/v2` v2.0.0](https://github.com/faustbrian/go-ecma-regexp/releases/tag/v2.0.0).
The original v1 module has no backport of the replacement admission repair;
see [security guidance](docs/security.md) for its advisory and migration.
Support windows are documented per module and in
[`COMPATIBILITY.md`](COMPATIBILITY.md).

## Security Gates

Releases require isolated tests, race and hostile-input checks, exact coverage
and mutation results, `govulncheck`, secret scanning, license verification,
SBOM generation, provenance validation, and clean-consumer resolution. A
missing scanner or unavailable service is a failed gate, not a warning.

Security fixes MUST include a regression test that does not publish weaponized
details or real secrets. Credentials MUST be redacted from logs and evidence.

## Repository Assurance

The repository [safety and concurrency policy](AGENTS.md#safety-and-concurrency)
and [supply-chain policy](AGENTS.md#dependencies-and-supply-chain) define shared
trust boundaries and release requirements. Package-specific security guidance
refines those rules for its owned boundary.
