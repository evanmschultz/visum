#!/usr/bin/env sh
set -eu

GOROOT="$(go env GOROOT)"
cp "$GOROOT/lib/wasm/wasm_exec.js" "$(dirname "$0")/../web/wasm_exec.js"

printf "Copied wasm_exec.js from %s\n" "$GOROOT"
