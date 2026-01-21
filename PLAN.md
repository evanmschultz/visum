# Visum Plan

## Goal
Build a Go + WebAssembly web app that renders times-table circles with full user control over parameters, animation, and stepwise navigation. The app must be extensible, maintainable, testable, and use hexagonal architecture where appropriate.

## Scope
- WASM renderer draws to an HTML canvas.
- Full control over points, multiplier, rotation, line count, colors, and labels.
- Animation controls for multiple variables with speed, loop, and ping-pong.
- Step controls for forward/backward navigation.
- Clean Go architecture with clear boundaries and tests for core math.

## Architecture
- `internal/core`: Pure math and geometry; no IO or JS dependencies.
- `internal/app`: Engine that owns state, applies updates, and produces frames.
- `internal/adapter/web`: WASM DOM + Canvas adapter, UI binding, and rendering.
- `cmd/visum`: WASM entrypoint.
- `cmd/visum-serve`: Local static file server for development.

## Phases
1. **Foundation**
   - Create module, directory structure, and docs.
   - Define domain types and frame model.

2. **Core + App Layer**
   - Implement math for points and lines.
   - Implement engine state, animations, and step controls.
   - Add unit tests for geometry and line generation.

3. **Web Adapter + UI**
   - Build HTML/CSS UI and a Go WASM adapter.
   - Wire inputs to engine state and render to canvas.
   - Provide requestAnimationFrame loop and resize handling.

4. **Polish**
   - Document setup, build, and usage.
   - Add MIT license and contribution guidelines.
   - Ensure extensible structure for future features.

## Deliverables
- `PLAN.md` (this file)
- `AGENTS.md`
- `README.md`
- `LICENSE` (MIT)
- `CONTRIBUTING.md`
- Fully working Go + WASM project with tests
