# Visum

Visum is a Go + WebAssembly app for exploring times-table circles (modular multiplication on a circle). It renders directly to an HTML canvas and exposes full, live controls for geometry, appearance, and animation.

## Features
- Go-only core logic with clean separation between math, state, and web adapters.
- Live control of point count, multiplier, rotation, line count, and styling.
- Multiple animation tracks (lines, multiplier, points) with speed, loop, and ping-pong.
- Stepwise forward/back control over lines, multiplier, or points.
- Responsive canvas with HiDPI support.

## Getting Started

### Prerequisites
- Go 1.22+

### Setup
1. Copy the Go WASM runtime shim:

```
./scripts/copy_wasm_exec.sh
```

2. Build the WASM binary:

```
GOOS=js GOARCH=wasm go build -o web/app.wasm ./cmd/visum
```

3. Start the dev server:

```
go run ./cmd/visum-serve
```

4. Open the app in your browser at `http://localhost:8080`.

## Usage

## Screenshots

_Placeholders (drop images in later):_

![Screenshot - Cardioid](assets/screens/cardioid.png)
![Screenshot - High multiplier](assets/screens/high-multiplier.png)
![Screenshot - Haeckel palette](assets/screens/haeckel-palette.png)

### Geometry
- **Points (N)**: Number of points around the circle.
- **Multiplier (k)**: Multiplies each index before mapping back to the circle.
- **Rotation**: Rotates the entire circle (degrees).
- **Start index**: Offset for line drawing.
- **Line count**: Draw only the first N lines for incremental builds.

### Appearance
- Toggle circle, points, and labels.
- Adjust line width and point radius.
- Customize colors for background, lines, circle, points, and labels.

### Animation
- Enable individual animations for lines, multiplier, and points.
- Each animation has a start, end, speed, and optional loop/ping-pong.
- Use play/pause plus step controls to move forward or backward.

### Step Controls
- Choose the target (lines, multiplier, points).
- Set the step amount.
- Use Step + / Step - to move manually, even while paused.

## Mathematical Notes
This visualization maps each point `n` on the circle to `k * n (mod N)`. The modulo keeps the mapping on the circle and produces repeating symmetries tied to the arithmetic structure of `N` and `k`.

The cardioid is the envelope of chords obtained by connecting each point to its double. Other multipliers yield related envelopes (such as nephroids), and fractional multipliers create intricate woven curves.

### Further Reading
- Cardioid times-table chords: https://divisbyzero.com/2018/04/02/i-heart-cardioids/
- Interactive modular circle explainer: https://www.redblobgames.com/x/1847-mathologer-modulo-circle/
- Modular arithmetic overview: https://en.wikipedia.org/wiki/Modular_arithmetic
- Cardioid geometry: https://en.wikipedia.org/wiki/Cardioid

## Architecture
This repo uses a hexagonal architecture:

- `internal/core`: Pure geometry and times-table math. No DOM, IO, or WebAssembly.
- `internal/app`: The engine with state, animations, and frame creation.
- `internal/adapter/web`: WASM adapter that binds DOM events and renders to canvas.
- `cmd/visum`: WASM entrypoint.
- `cmd/visum-serve`: Local static server.

## Testing

```
go test ./internal/...
```

## Development Notes
- `web/wasm_exec.js` must match your local Go version.
- `web/app.wasm` is generated and should not be edited by hand.
- Canvas sizing uses device pixel ratio for crisp results.

## Roadmap Ideas
- Export PNG/SVG snapshots.
- Add presets and saved configurations.
- Support animated color palettes.

## License
MIT. See `LICENSE`.

## Contributing
See `CONTRIBUTING.md` for workflow, testing, and style expectations.
