#!/usr/bin/env bash
set -euo pipefail

# The shared runner supplies V8 runtimes but has no JavaScriptCore bootstrap.
# Provision immutable test-only peers at the existing interoperability boundary.
node_version=24.4.1
bun_version=1.3.14
case "$(uname -s)-$(uname -m)" in
  Darwin-arm64)
    node_platform=darwin-arm64
    node_sha=55a772a600b7bdafb4b35945b3935090e27aff9934b4c11b281220fcd99139d7
    bun_platform=darwin-aarch64
    bun_sha=d8b96221828ad6f97ac7ac0ab7e95872341af763001e8803e8267652c2652620
    ;;
  Darwin-x86_64)
    node_platform=darwin-x64
    node_sha=59fbad953a0705e78d220079fb6d10d341d0a61afd3aeb4db2a87207fddd8944
    bun_platform=darwin-x64
    bun_sha=4183df3374623e5bab315c547cfa0974533cd457d86b73b639f7a87974cd6633
    ;;
  Linux-aarch64|Linux-arm64)
    node_platform=linux-arm64
    node_sha=fde5421e2652e51199bc678e1e6c4d80bbb4c55337ec0a82206568517e9792ef
    bun_platform=linux-aarch64
    bun_sha=a27ffb63a8310375836e0d6f668ae17fa8d8d18b88c37c821c65331973a19a3b
    ;;
  Linux-x86_64)
    node_platform=linux-x64
    node_sha=063f2eb299ba60e3fc9b424d8e87d0e2f6be84b39bdeadc421ee2865914c498b
    bun_platform=linux-x64
    bun_sha=951ee2aee855f08595aeec6225226a298d3fea83a3dcd6465c09cbccdf7e848f
    ;;
  *) echo "unsupported differential provisioning platform" >&2; exit 1 ;;
esac

task="$(mktemp -d "${TMPDIR:-/tmp}/ecma-differential.XXXXXX")"
cleanup() {
  local status=$?
  trap - EXIT HUP INT TERM
  chmod -R u+w "${task}"
  find "${task}" -depth -delete
  exit "${status}"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' HUP TERM
mkdir -p "${task}/bin" "${task}/build" "${task}/mod" "${task}/tmp"

verify() {
  local archive="$1" expected="$2" actual
  actual="$(shasum -a 256 "${archive}")"
  [[ "${actual%% *}" == "${expected}" ]] || {
    echo "differential peer archive checksum mismatch" >&2
    return 1
  }
}
curl --fail --silent --show-error --location --connect-timeout 10 --max-time 300 --retry 3 \
  "https://nodejs.org/dist/v${node_version}/node-v${node_version}-${node_platform}.tar.gz" \
  --output "${task}/node.tar.gz"
verify "${task}/node.tar.gz" "${node_sha}"
tar -xzf "${task}/node.tar.gz" -C "${task}/bin" --strip-components=2 \
  "node-v${node_version}-${node_platform}/bin/node"
curl --fail --silent --show-error --location --connect-timeout 10 --max-time 300 --retry 3 \
  "https://github.com/oven-sh/bun/releases/download/bun-v${bun_version}/bun-${bun_platform}.zip" \
  --output "${task}/bun.zip"
verify "${task}/bun.zip" "${bun_sha}"
unzip -q -j "${task}/bun.zip" "bun-${bun_platform}/bun" -d "${task}/bin"
export PATH="${task}/bin:${PATH}"
[[ "$(node --version)" == "v${node_version}" ]]
[[ "$(bun --version)" == "${bun_version}" ]]
node --version
bun --version
if command -v deno >/dev/null; then
  [[ "$(deno --version | head -n 1)" == 'deno 2.9.3' ]]
  deno --version
fi

export GOCACHE="${task}/build" GOMODCACHE="${task}/mod" GOTMPDIR="${task}/tmp" GOWORK=off
export ECMA_RELEASE_DIFFERENTIAL=1
go test -mod=readonly . -run '^(TestDifferentialMatchingAgainstJavaScriptEngines|TestOverlappingLibraryDifferential|TestReleaseDifferentialRejectsTwoV8Runtimes)$' -count=1 -timeout=15m -v
