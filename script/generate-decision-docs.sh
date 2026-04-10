#!/bin/sh

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
META_DIR=${1:-"$ROOT_DIR/doc/decision-meta"}
DOC_DIR=${2:-"$ROOT_DIR/doc/decision"}
SEER_BIN=${SEER_BIN:-"$ROOT_DIR/.e2e-bin/seer"}

if [ ! -x "$SEER_BIN" ]; then
	printf 'error: seer binary not found or not executable at %s\n' "$SEER_BIN" >&2
	printf 'hint: run make build-go or set SEER_BIN to an executable seer binary\n' >&2
	exit 1
fi

mkdir -p "$META_DIR"
mkdir -p "$DOC_DIR"

decision_files=$(find "$META_DIR" -type f -name '*.seer.cue' | LC_ALL=C sort)

if [ -z "$decision_files" ]; then
	printf 'decision_configs=0\n'
	exit 0
fi

printf '%s\n' "$decision_files" | while IFS= read -r config_path; do
	relative_path=${config_path#"$META_DIR"/}
	output_path="$DOC_DIR/${relative_path%.seer.cue}.md"
	output_dir=$(dirname "$output_path")

	mkdir -p "$output_dir"

	"$SEER_BIN" validate --config "$config_path"
	"$SEER_BIN" report generate --config "$config_path" > "$output_path"

	printf 'generated=%s\n' "$output_path"
done
