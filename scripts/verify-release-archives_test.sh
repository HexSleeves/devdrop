#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
verify="$repo_root/scripts/verify-release-archives.sh"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

snapshot_plan=$(make -s -n -C "$repo_root" snapshot)
build_line=$(grep -nF './build-all.sh' <<<"$snapshot_plan" | cut -d: -f1)
release_line=$(grep -nF 'goreleaser release --snapshot --clean --skip=publish' <<<"$snapshot_plan" | cut -d: -f1)
if [[ -z "$build_line" || -z "$release_line" || "$build_line" -ge "$release_line" ]]; then
  echo "snapshot must build companion binaries before GoReleaser" >&2
  exit 1
fi

make_archive() {
  local os=$1 arch=$2 include_tui=${3:-1} root="$tmp/root"
  rm -rf "$root"
  mkdir -p "$root"
  printf '#!/bin/sh\n' >"$root/devspace"
  chmod 755 "$root/devspace"
  if [[ "$include_tui" == 1 ]]; then
    printf '#!/bin/sh\n' >"$root/devspace-tui"
    chmod 755 "$root/devspace-tui"
  fi
  tar -C "$root" -czf "$tmp/devspace_v0.0.0_${os}_${arch}.tar.gz" .
}

for target in linux_amd64 linux_arm64 darwin_amd64 darwin_arm64; do
  make_archive "${target%_*}" "${target#*_}"
done
"$verify" "$tmp"

make_archive darwin arm64 0
if "$verify" "$tmp" >/dev/null 2>&1; then
  echo "expected missing companion validation to fail" >&2
  exit 1
fi
