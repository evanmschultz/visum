set shell := ["/bin/sh", "-c"]

# Default task
_default:
	@just --list

# Copy the Go WASM runtime shim
wasm-shim:
	./scripts/copy_wasm_exec.sh

# Build the WASM binary
wasm-build: wasm-shim
	GOCACHE="$PWD/.go-cache" GOOS=js GOARCH=wasm go build -o web/app.wasm ./cmd/visum

# Run the local dev server
serve addr=":8080":
	GOCACHE="$PWD/.go-cache" go run ./cmd/visum-serve --addr {{addr}}

# Build then serve
serve-wasm addr=":8080": wasm-build
	GOCACHE="$PWD/.go-cache" go run ./cmd/visum-serve --addr {{addr}}

# Run unit tests (native)
test:
	GOCACHE="$PWD/.go-cache" go test ./...

# Run wasm tests
wasm-test:
	GOCACHE="$PWD/.go-cache" GOOS=js GOARCH=wasm go test -exec="$(go env GOROOT)/lib/wasm/go_js_wasm_exec" ./internal/adapter/web

# Coverage (native)
cover:
	GOCACHE="$PWD/.go-cache" GOTOOLCHAIN=go1.25.6+auto go test -cover ./...

# Coverage (wasm)
wasm-cover:
	GOCACHE="$PWD/.go-cache" GOTOOLCHAIN=go1.25.6+auto GOOS=js GOARCH=wasm go test -cover -exec="$(go env GOROOT)/lib/wasm/go_js_wasm_exec" ./internal/adapter/web
