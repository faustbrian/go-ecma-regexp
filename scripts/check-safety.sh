#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
task="$(mktemp -d "${TMPDIR:-/tmp}/ecma-regexp-safety.XXXXXX")"
cleanup() {
	find "${task}" -depth -delete 2>/dev/null || true
}
trap cleanup EXIT HUP INT TERM

files="${task}/go-files"
if ! find "${root}" \
	\( \
		-path "${root}/.git" -o \
		-path "${root}/.golib-tooling" -o \
		-path "${root}/.verification" -o \
		-name '.*' -o \
		-name '_*' -o \
		-name 'testdata' -o \
		-name 'vendor' \
	\) -prune -o \
	-type f -name '*.go' ! -name '*_test.go' -print0 >"${files}"; then
	echo "failed to inventory production Go files" >&2
	exit 1
fi

production_files="${task}/production-go-files"
if ! : >"${production_files}"; then
	echo "failed to prepare production Go file inventory" >&2
	exit 1
fi
while IFS= read -r -d '' file; do
	directory="$(dirname "${file}")"
	nested=0
	while [[ "${directory}" != "${root}" && "${directory}" != / ]]; do
		if [[ -f "${directory}/go.mod" ]]; then
			nested=1
			break
		fi
		directory="$(dirname "${directory}")"
	done
	if (( nested == 0 )); then
		printf '%s\0' "${file}" >>"${production_files}"
	fi
done <"${files}"

if [[ ! -s "${production_files}" ]]; then
	echo "production Go file inventory is empty" >&2
	exit 1
fi

scan() {
	local pattern="$1"
	local status
	while IFS= read -r -d '' file; do
		set +e
		grep -nE "${pattern}" "${file}"
		status=$?
		set -e
		case "${status}" in
			0)
				echo "forbidden production mechanism found" >&2
				return 1
				;;
			1)
				;;
			*)
				echo "production safety scan failed closed" >&2
				return "${status}"
				;;
		esac
	done <"${production_files}"
}

scan '(^|[^[:alnum:]_])(unsafe|C)\.|go:linkname|func[[:space:]]+init[[:space:]]*\('
scan '^[[:space:]]*go[[:space:]]+'
