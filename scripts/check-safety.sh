#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
task="$(mktemp -d "${TMPDIR:-/tmp}/ecma-regexp-safety.XXXXXX")"
cleanup() {
	find "${task}" -depth -delete 2>/dev/null || true
}
trap cleanup EXIT HUP INT TERM

nested_modules="${task}/nested-modules"
if ! find "${root}" -mindepth 2 -type f -name go.mod -print >"${nested_modules}"; then
	echo "failed to inventory nested modules for production safety scan" >&2
	exit 1
fi

scan_args=(
	--hidden
	--no-ignore
	--glob '*.go'
	--glob '!**/*_test.go'
	--glob '!**/.*'
	--glob '!**/_*'
	--glob '!**/_*/**'
	--glob '!**/testdata/**'
	--glob '!**/vendor/**'
	--glob '!**/.golib-tooling/**'
	--glob '!**/.verification/**'
)
while IFS= read -r module; do
	[[ -n "${module}" ]] || continue
	relative="${module#"${root}"/}"
	scan_args+=(--glob "!${relative%/go.mod}/**")
done <"${nested_modules}"

scan() {
	local pattern="$1"
	local status
	set +e
	rg -n "${pattern}" "${scan_args[@]}" "${root}"
	status=$?
	set -e
	case "${status}" in
		0)
			echo "forbidden production mechanism found" >&2
			return 1
			;;
		1)
			return 0
			;;
		*)
			echo "production safety scan failed closed" >&2
			return "${status}"
			;;
	esac
}

scan '(^|[^[:alnum:]_])(unsafe|C)\.|go:linkname|func init\s*\('
scan '^[[:space:]]*go[[:space:]]+'
