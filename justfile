set shell := ["/bin/sh", "-c"]
set env := {"GOCACHE": "{{justfile_directory()}}/.go-cache"}

# Default task
_default:
	@just --list

# Copy the Go WASM runtime shim
wasm-shim:
	./scripts/copy_wasm_exec.sh

# Build the WASM binary
wasm-build: wasm-shim
	GOOS=js GOARCH=wasm go build -o web/app.wasm ./cmd/visum

# Run the local dev server
serve addr=":8080":
	go run ./cmd/visum-serve --addr {{addr}}

# Build then serve
serve-wasm addr=":8080": wasm-build
	go run ./cmd/visum-serve --addr {{addr}}

# Run unit tests (native)
test:
	go test ./...

# Run wasm tests
wasm-test:
	GOOS=js GOARCH=wasm go test -exec="$(go env GOROOT)/lib/wasm/go_js_wasm_exec" ./internal/adapter/web

# Coverage (native)
cover:
	GOTOOLCHAIN=go1.25.6+auto go test -cover ./...

# Coverage (wasm)
wasm-cover:
	GOTOOLCHAIN=go1.25.6+auto GOOS=js GOARCH=wasm go test -cover -exec="$(go env GOROOT)/lib/wasm/go_js_wasm_exec" ./internal/adapter/web
