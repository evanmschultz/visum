//go:build js && wasm

package web

import (
	"errors"
	"fmt"
	"math"
	"syscall/js"

	"github.com/evanschultz/visum/internal/core"
)

// CanvasRenderer draws frames onto an HTML canvas.
type CanvasRenderer struct {
	canvas js.Value
	ctx    js.Value
	cssSize core.Size
	fonts   string
}

// NewCanvasRenderer locates the canvas by ID and prepares a 2D context.
func NewCanvasRenderer(canvasID string) (*CanvasRenderer, error) {
	doc := js.Global().Get("document")
	if doc.IsUndefined() {
		return nil, errors.New("document not available")
	}
	canvas := doc.Call("getElementById", canvasID)
	if canvas.IsNull() || canvas.IsUndefined() {
		return nil, fmt.Errorf("canvas %q not found", canvasID)
	}
	ctx := canvas.Call("getContext", "2d")
	if ctx.IsNull() || ctx.IsUndefined() {
		return nil, errors.New("2d context not available")
	}

	return &CanvasRenderer{
		canvas: canvas,
		ctx:    ctx,
		fonts:  "300 12px \"Source Serif 4\", \"Iowan Old Style\", \"Palatino Linotype\", serif",
	}, nil
}

// Size returns the current canvas size in CSS pixels.
func (r *CanvasRenderer) Size() core.Size {
	return r.cssSize
}

// EnsureSize syncs the canvas backing store with its CSS size.
func (r *CanvasRenderer) EnsureSize() {
	rect := r.canvas.Call("getBoundingClientRect")
	cssWidth := rect.Get("width").Float()
	cssHeight := rect.Get("height").Float()
	if cssWidth <= 0 || cssHeight <= 0 {
		return
	}
	dpr := js.Global().Get("devicePixelRatio").Float()
	if dpr <= 0 {
		dpr = 1
	}

	pixelWidth := math.Round(cssWidth * dpr)
	pixelHeight := math.Round(cssHeight * dpr)
	if r.canvas.Get("width").Float() != pixelWidth {
		r.canvas.Set("width", pixelWidth)
	}
	if r.canvas.Get("height").Float() != pixelHeight {
		r.canvas.Set("height", pixelHeight)
	}

	r.ctx.Call("setTransform", dpr, 0, 0, dpr, 0, 0)
	r.cssSize = core.Size{Width: cssWidth, Height: cssHeight}
}

// Render draws the frame using the provided params for styling.
func (r *CanvasRenderer) Render(frame core.Frame, params core.Params) {
	r.EnsureSize()
	if r.cssSize.Width <= 0 || r.cssSize.Height <= 0 {
		return
	}

	ctx := r.ctx
	ctx.Set("fillStyle", params.Colors.Background)
	ctx.Call("fillRect", 0, 0, r.cssSize.Width, r.cssSize.Height)

	ctx.Set("lineWidth", params.LineWidth)
	ctx.Set("lineCap", "round")
	ctx.Set("strokeStyle", params.Colors.Line)

	for _, line := range frame.Lines {
		ctx.Call("beginPath")
		ctx.Call("moveTo", line.From.X, line.From.Y)
		ctx.Call("lineTo", line.To.X, line.To.Y)
		ctx.Call("stroke")
	}

	if params.ShowCircle {
		ctx.Set("strokeStyle", params.Colors.Circle)
		ctx.Call("beginPath")
		ctx.Call("arc", frame.Circle.Center.X, frame.Circle.Center.Y, frame.Circle.Radius, 0, 2*math.Pi)
		ctx.Call("stroke")
	}

	if params.ShowPoints {
		ctx.Set("fillStyle", params.Colors.Point)
		for _, point := range frame.Points {
			ctx.Call("beginPath")
			ctx.Call("arc", point.X, point.Y, params.PointRadius, 0, 2*math.Pi)
			ctx.Call("fill")
		}
	}

	if params.ShowLabels {
		fontSize := math.Max(10, frame.Circle.Radius*0.06)
		r.fonts = fmt.Sprintf("300 %.0fpx \"Source Serif 4\", \"Iowan Old Style\", \"Palatino Linotype\", serif", fontSize)
		ctx.Set("font", r.fonts)
		ctx.Set("fillStyle", params.Colors.Label)
		ctx.Set("textAlign", "center")
		ctx.Set("textBaseline", "middle")
		for _, label := range frame.Labels {
			ctx.Call("fillText", label.Text, label.Position.X, label.Position.Y)
		}
	}
}
