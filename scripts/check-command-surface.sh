#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

# This is deliberately an allowlist. Completed SDD specs/proofs/validations,
# docs/architecture/manifest-merge.md (a superseded spike), generated demo
# GIFs, and the intentional negative fixtures under scripts/testdata/ are not
# scanned. The linked capstone HTML remains part of the maintained surface.
maintained_files=(
  AGENTS.md
  CLAUDE.md
  Makefile
  README.md
  ARCHITECTURE.md
  docs/architecture/access-roles.md
  docs/architecture/fuse-lazy-mount.md
  docs/operations/macos-fuse-run-playbook.md
  docs/operations/release-readiness.md
  docs/operations/release.md
  docs/capstone/README.md
  docs/capstone/case-study.md
  docs/capstone/demo-script.md
  docs/capstone/index.html
  docs/capstone/playbook-contribution.md
  docs/capstone/proof-artifacts.md
  docs/capstone/remote-agent-case-study.md
  docs/capstone/spec.md
  docs/demos/README.md
  docs/demos/capstone-runbook.md
  docs/demos/capstone-rehearsal.sh
  docs/demos/hosted-sync-serve.sh
  docs/demos/mkgitproject.sh
  docs/demos/project-remote.sh
  docs/demos/set-api-key.sh
  docs/demos/spawn-project-later.sh
  docs/demos/capstone-walkthrough.tape
  docs/demos/env-secrets.tape
  docs/demos/getting-started.tape
  docs/demos/hosted-sync.tape
  docs/demos/mount-preview.tape
  docs/demos/project-lifecycle.tape
  docs/demos/reconcile.tape
  docs/demos/setup-commands.tape
  docs/demos/ui-dashboard.tape
  docs/demos/watch.tape
  docs/demos/workspace-sync.tape
  scripts/demo-check.sh
)

removed_patterns=(
  'workspace[[:space:]]+(push|pull|diff|reconcile|remote|scan|sync)([^[:alnum:]_-]|$)'
  'project[[:space:]]+(add|remove|hydrate)([^[:alnum:]_-]|$)'
  '(^|[^[:alnum:]_-])(devspace|bin/devspace|\./bin/devspace|\$DS|"\$devspace")[[:space:]]+project[[:space:]]+status([^[:alnum:]_-]|$)'
  'env[[:space:]]+pull([^[:alnum:]_-]|$)'
  'setup[[:space:]]+(plan|apply)([^[:alnum:]_-]|$)'
  '(^|[^[:alnum:]_-])(devspace|bin/devspace|\./bin/devspace|\$DS|"\$devspace")[[:space:]]+hosted[[:space:]]+serve([^[:alnum:]_-]|$)'
  '(^|[^[:alnum:]_-])(devspace|bin/devspace|\./bin/devspace|\$DS|"\$devspace")[[:space:]]+mount([^[:alnum:]_-]|$)'
  '(^|[^[:alnum:]_-])devspace[[:space:]]+tui([^[:alnum:]_-]|$)'
  '(^|[^[:alnum:]_-])devspace[[:space:]]+version([^[:alnum:]_-]|$)'
  '(^|[^[:alnum:]_-])devspace[[:space:]]+workspace([[:space:]]+--json)?[[:space:]]*[`]'
  '(^|[^[:alnum:]_-])devspace[[:space:]]+project([[:space:]]+--json)?[[:space:]]*[`]'
  '(^|[^[:alnum:]_-])devspace[[:space:]]+workspace([[:space:]]+--json)?[[:space:]`"'"'"']*$'
  '(^|[^[:alnum:]_-])devspace[[:space:]]+project([[:space:]]+--json)?[[:space:]`"'"'"']*$'
)

canonical_patterns=(
  'sync[[:space:]]+remote'
  'sync[[:space:]]+push'
  'sync[[:space:]]+pull'
  'sync[[:space:]]+diff'
  'sync[[:space:]]+reconcile'
  'project[[:space:]]+list'
  'project[[:space:]]+track'
  'project[[:space:]]+untrack'
  'project[[:space:]]+update'
  'env[[:space:]]+write'
  'setup[[:space:]]+show'
  'setup[[:space:]]+run'
  'experimental[[:space:]]+hosted[[:space:]]+serve'
  'experimental[[:space:]]+mount'
  'devspace[[:space:]]+ui'
  'devspace[[:space:]]+--version'
  'status[[:space:]]+client-a-api'
)

# README's marked pre-1.0 migration table names removed paths as historical
# labels. Only that bounded block is omitted from command matching.
scan_file() {
  local exclude_markers=0
  [[ "$1" == "README.md" || "$1" == "docs/capstone/index.html" ]] && exclude_markers=1
  awk -v exclude_markers="$exclude_markers" '
    exclude_markers && /command-surface-migration:start/ { excluded = 1; next }
    exclude_markers && /command-surface-migration:end/ { excluded = 0; next }
    !excluded { print FNR ":" $0 }
  ' "$1"
}

# Match against one whitespace-normalized stream so Markdown line wrapping
# cannot split a command path into separately clean physical lines.
normalized_file() {
  local exclude_markers=0
  [[ "$1" == "README.md" || "$1" == "docs/capstone/index.html" ]] && exclude_markers=1
  awk -v exclude_markers="$exclude_markers" '
    exclude_markers && /command-surface-migration:start/ { excluded = 1; next }
    exclude_markers && /command-surface-migration:end/ { excluded = 0; next }
    !excluded {
      line = $0
      gsub(/[[:space:]]+/, " ", line)
      sub(/^ /, "", line)
      sub(/ $/, "", line)
      if (length(line)) printf "%s ", line
    }
    END { print "" }
  ' "$1"
}

check_removed_paths() {
  local found=0 pattern file matches normalized_match
  for pattern in "${removed_patterns[@]}"; do
    for file in "$@"; do
      matches="$(scan_file "$file" | grep -E "$pattern" || true)"
      normalized_match=0
      if normalized_file "$file" | grep -E "$pattern" >/dev/null; then
        normalized_match=1
      fi
      if [[ -n "$matches" ]] || ((normalized_match)); then
        if [[ -z "$matches" ]]; then
          echo "command-surface: removed path (wrapped): $file: pattern $pattern" >&2
        fi
        while IFS= read -r match; do
          [[ -z "$match" ]] || echo "command-surface: removed path: $file:$match" >&2
        done <<<"$matches"
        found=1
      fi
    done
  done
  ((found == 0))
}

case "${1:-}" in
  --self-test)
    for fixture in \
      scripts/testdata/removed-command-wrapped.md \
      scripts/testdata/removed-command-bare-workspace.md \
      scripts/testdata/removed-command-bare-project.md \
      scripts/testdata/removed-command-inline-workspace.md \
      scripts/testdata/removed-command-inline-project.md \
      scripts/testdata/removed-command-fake-markers.md; do
      if check_removed_paths "$fixture" >/dev/null 2>&1; then
        echo "command-surface self-test: removed path was not rejected: $fixture" >&2
        exit 1
      fi
    done
    if ! check_removed_paths scripts/testdata/canonical-command-inline.md >/dev/null 2>&1; then
      echo "command-surface self-test: canonical inline path was rejected" >&2
      exit 1
    fi
    echo "command-surface self-test: wrapped, bare, and inline removed paths rejected"
    exit 0
    ;;
  "") ;;
  *)
    echo "usage: scripts/check-command-surface.sh [--self-test]" >&2
    exit 2
    ;;
esac

failed=0
for file in "${maintained_files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "command-surface: missing maintained file: $file" >&2
    failed=1
  fi
done

if ((failed)); then
  exit 1
fi

if [[ "$(grep -c '^<!-- command-surface-migration:start -->$' README.md)" -ne 1 ||
      "$(grep -c '^<!-- command-surface-migration:end -->$' README.md)" -ne 1 ]]; then
  echo "command-surface: README migration exclusion markers must be one start/end pair" >&2
  exit 1
fi

if [[ "$(grep -c 'command-surface-migration:start' docs/capstone/index.html)" -ne 3 ||
      "$(grep -c 'command-surface-migration:end' docs/capstone/index.html)" -ne 3 ]]; then
  echo "command-surface: capstone HTML migration exclusions must be three generated copies" >&2
  exit 1
fi

if ! check_removed_paths "${maintained_files[@]}"; then
  failed=1
fi

for pattern in "${canonical_patterns[@]}"; do
  found=0
  for file in "${maintained_files[@]}"; do
    if normalized_file "$file" | grep -E "$pattern" >/dev/null; then
      found=1
      break
    fi
  done
  if (( ! found )); then
    echo "command-surface: missing canonical path matching: $pattern" >&2
    failed=1
  fi
done

if ((failed)); then
  exit 1
fi

echo "command-surface: maintained documentation and demos use canonical commands"
