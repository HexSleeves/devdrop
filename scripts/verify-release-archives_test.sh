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

mkdir -p "$tmp/bin"
printf 'package main\nfunc main() {}\n' >"$tmp/main.go"
for target in linux_amd64 linux_arm64 darwin_amd64 darwin_arm64; do
  CGO_ENABLED=0 GOOS=${target%_*} GOARCH=${target#*_} \
    go build -trimpath -o "$tmp/bin/$target" "$tmp/main.go"
done

make_archive() {
  local os=$1 arch=$2 include_tui=${3:-1}
  local primary_target=${4:-${os}_${arch}} tui_target=${5:-${os}_${arch}}
  local primary_mode=${6:-755} version=${7:-v0.0.0} root="$tmp/root"
  rm -rf "$root"
  mkdir -p "$root"
  cp "$tmp/bin/$primary_target" "$root/devspace"
  chmod "$primary_mode" "$root/devspace"
  if [[ $include_tui == 1 ]]; then
    cp "$tmp/bin/$tui_target" "$root/devspace-tui"
    chmod 755 "$root/devspace-tui"
  fi
  tar -C "$root" -czf "$tmp/devspace_${version}_${os}_${arch}.tar.gz" .
}

write_checksums() {
  (
    cd "$tmp"
    if command -v sha256sum >/dev/null 2>&1; then
      sha256sum devspace_*.tar.gz >checksums.txt
    else
      shasum -a 256 devspace_*.tar.gz >checksums.txt
    fi
  )
}

expect_failure() {
  local label=$1 expected=$2 output
  if output=$("$verify" "$tmp" 2>&1); then
    echo "expected $label validation to fail" >&2
    exit 1
  elif [[ $output != *"$expected"* ]]; then
    echo "unexpected $label error: $output" >&2
    exit 1
  fi
}

for target in linux_amd64 linux_arm64 darwin_amd64 darwin_arm64; do
  make_archive "${target%_*}" "${target#*_}"
done
write_checksums
"$verify" "$tmp"

make_archive linux amd64 1 linux_arm64 linux_amd64
write_checksums
expect_failure "wrong-architecture primary" "devspace: expected linux/amd64"

make_archive linux amd64 1 darwin_amd64 linux_amd64
write_checksums
expect_failure "wrong-OS primary" "devspace: expected linux/amd64"

make_archive linux amd64 1 linux_amd64 linux_arm64
write_checksums
expect_failure "wrong-architecture companion" "devspace-tui: expected linux/amd64"

make_archive linux amd64
make_archive darwin arm64 1 darwin_arm64 linux_arm64
write_checksums
expect_failure "wrong-OS companion" "devspace-tui: expected darwin/arm64"

make_archive linux amd64
make_archive darwin arm64 0
write_checksums
expect_failure "missing companion" "missing executable devspace-tui"

make_archive darwin arm64 1 darwin_arm64 darwin_arm64 644
write_checksums
expect_failure "non-executable primary" "missing executable devspace"

make_archive darwin arm64
write_checksums
sed -n '2,4p' "$tmp/checksums.txt" >"$tmp/checksums.txt.clean"
mv "$tmp/checksums.txt.clean" "$tmp/checksums.txt"
expect_failure "missing checksum" "expected four archive entries, found 3"

write_checksums
{
  sed -n '1,3p' "$tmp/checksums.txt"
  sed -n '1p' "$tmp/checksums.txt"
} >"$tmp/checksums.txt.clean"
mv "$tmp/checksums.txt.clean" "$tmp/checksums.txt"
expect_failure "duplicate checksum" "missing or invalid checksum"

write_checksums
printf '%064d  devspace-tui_darwin_arm64\n' 0 >>"$tmp/checksums.txt"
expect_failure "extra checksum entry" "expected four archive entries, found 5"

write_checksums
{
  printf '%064d  %s\n' 0 "$(awk 'NR == 1 { print $2 }' "$tmp/checksums.txt")"
  sed -n '2,4p' "$tmp/checksums.txt"
} >"$tmp/checksums.txt.clean"
mv "$tmp/checksums.txt.clean" "$tmp/checksums.txt"
expect_failure "checksum mismatch" "missing or invalid checksum"

mv "$tmp/devspace_v0.0.0_darwin_arm64.tar.gz" "$tmp/devspace_v9.9.9_darwin_arm64.tar.gz"
write_checksums
expect_failure "inconsistent archive version" "inconsistent archive version"
