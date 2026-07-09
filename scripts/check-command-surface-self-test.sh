#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
output="$("$repo_root/scripts/check-command-surface.sh" --self-test)"

[[ "$output" == "command-surface self-test: wrapped removed path rejected" ]] || {
  echo "command-surface self-test failed: $output" >&2
  exit 1
}

echo "$output"
