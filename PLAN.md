# Plan: Client-Side Real-Time Video Export

## Goals
- Keep all still exports (PNG/WebP/SVG) client-side.
- Provide **real-time** video export using `MediaRecorder`.
- Add clear UX warnings to keep the tab open during export.
- Maintain hexagonal architecture and ~80%+ test coverage.

## Scope
1) **Image exports (client)**
   - PNG/WebP from the live canvas.
   - SVG from Go geometry via WASM bridge.
2) **Video exports (client)**
   - **Export video (real-time):** timed clip based on animation bounds.
   - **Record video (manual):** live capture until stop.

## UX Plan
- Progress bar + ETA for timed exports.
- Indeterminate progress + elapsed timer for manual recordings.
- Warning modal on any interaction while recording.
- `beforeunload` warning to prevent accidental refresh/close.

## Implementation Steps
1) Add progress/ETA UI and hover explainers.
2) Add warning modal and `beforeunload` handler.
3) Use animation bounds to calculate total duration for timed export.
4) Keep tests passing; update docs.

## Tests
- Keep existing Go + WASM tests passing.
- Maintain coverage targets (~80%).
