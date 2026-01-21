# AGENTS.md

## Project Overview
Visum is a Go + WebAssembly app that renders times-table circles in the browser with live, customizable controls and animation. The project is organized using a hexagonal architecture: core math and state live in the center, and the web adapter is a thin shell.

## Structure
- `internal/core`: Geometry and times-table math (pure functions).
- `internal/app`: Engine that owns state, animations, and frame generation.
- `internal/adapter/web`: WASM adapter for DOM inputs and Canvas rendering.
- `cmd/visum`: WASM entrypoint.
- `cmd/visum-serve`: Local static server for development.
- `web`: HTML/CSS/JS assets and the built WASM binary.

## Commands
- Build WASM: `GOOS=js GOARCH=wasm go build -o web/app.wasm ./cmd/visum`
- Run server: `go run ./cmd/visum-serve`
- Tests: `go test ./internal/...`

## Conventions
- Keep `internal/core` pure and deterministic (no JS or IO).
- Avoid leaking DOM-specific details into `internal/app`.
- Prefer small, focused functions and explicit state transitions.
- Update docs and tests alongside behavioral changes.

## Notes
- WASM runtime requires `wasm_exec.js` from the local Go toolchain.
- Canvas sizing uses `devicePixelRatio` for crisp rendering on HiDPI screens.
