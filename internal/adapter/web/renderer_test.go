//go:build js && wasm

package web

import (
	"testing"
	"syscall/js"

	"github.com/evanschultz/visum/internal/core"
)

func TestNewCanvasRendererMissingCanvas(t *testing.T) {
	setupDocument(t, map[string]js.Value{})
	if _, err := NewCanvasRenderer("missing"); err == nil {
		t.Fatalf("expected error when canvas is missing")
	}
}

func TestNewCanvasRendererMissingDocument(t *testing.T) {
	previous := js.Global().Get("document")
	js.Global().Set("document", js.Undefined())
	t.Cleanup(func() { js.Global().Set("document", previous) })

	if _, err := NewCanvasRenderer("visum-canvas"); err == nil {
		t.Fatalf("expected error when document is missing")
	}
}

func TestNewCanvasRendererMissingContext(t *testing.T) {
	canvas := js.ValueOf(map[string]interface{}{})
	getContext := js.FuncOf(func(this js.Value, args []js.Value) interface{} { return js.Null() })
	canvas.Set("getContext", getContext)
	setupDocument(t, map[string]js.Value{"visum-canvas": canvas})
	defer getContext.Release()

	if _, err := NewCanvasRenderer("visum-canvas"); err == nil {
		t.Fatalf("expected error when context is missing")
	}
}

func TestRenderWithStubCanvas(t *testing.T) {
	canvas := newStubCanvas(t)
	setupDocument(t, map[string]js.Value{"visum-canvas": canvas})
	js.Global().Set("devicePixelRatio", 1)

	renderer, err := NewCanvasRenderer("visum-canvas")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	params := core.DefaultParams()
	frame := core.BuildFrame(params, core.Size{Width: 800, Height: 600})
	renderer.Render(frame, params)

	if renderer.Size().Width != 800 || renderer.Size().Height != 600 {
		t.Fatalf("unexpected renderer size: %+v", renderer.Size())
	}
}

func TestRenderPointsAndLabels(t *testing.T) {
	canvas := newStubCanvas(t)
	setupDocument(t, map[string]js.Value{"visum-canvas": canvas})
	js.Global().Set("devicePixelRatio", 0)

	renderer, err := NewCanvasRenderer("visum-canvas")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	params := core.DefaultParams()
	params.ShowPoints = true
	params.ShowLabels = true
	params.LabelStep = 5
	frame := core.BuildFrame(params, core.Size{Width: 800, Height: 600})
	renderer.Render(frame, params)

	if renderer.Size().Width != 800 || renderer.Size().Height != 600 {
		t.Fatalf("unexpected renderer size after render: %+v", renderer.Size())
	}
}

func TestEnsureSizeZero(t *testing.T) {
	canvas := newZeroCanvas(t)
	setupDocument(t, map[string]js.Value{"visum-canvas": canvas})
	js.Global().Set("devicePixelRatio", 1)

	renderer, err := NewCanvasRenderer("visum-canvas")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	renderer.EnsureSize()
	if renderer.Size().Width != 0 || renderer.Size().Height != 0 {
		t.Fatalf("expected size to remain zero, got %+v", renderer.Size())
	}
}

func setupDocument(t *testing.T, elements map[string]js.Value) {
	doc := js.ValueOf(map[string]interface{}{})
	get := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return js.Null()
		}
		id := args[0].String()
		if el, ok := elements[id]; ok {
			return el
		}
		return js.Null()
	})
	doc.Set("getElementById", get)
	doc.Set("activeElement", js.Null())
	js.Global().Set("document", doc)

	t.Cleanup(func() {
		get.Release()
	})
}

func newStubCanvas(t *testing.T) js.Value {
	canvas := js.ValueOf(map[string]interface{}{})
	ctx := js.ValueOf(map[string]interface{}{})

	funcs := []js.Func{
		nopFunc(), // setTransform
		nopFunc(), // fillRect
		nopFunc(), // beginPath
		nopFunc(), // moveTo
		nopFunc(), // lineTo
		nopFunc(), // stroke
		nopFunc(), // arc
		nopFunc(), // fill
		nopFunc(), // fillText
	}

	ctx.Set("setTransform", funcs[0])
	ctx.Set("fillRect", funcs[1])
	ctx.Set("beginPath", funcs[2])
	ctx.Set("moveTo", funcs[3])
	ctx.Set("lineTo", funcs[4])
	ctx.Set("stroke", funcs[5])
	ctx.Set("arc", funcs[6])
	ctx.Set("fill", funcs[7])
	ctx.Set("fillText", funcs[8])

	getContext := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return ctx
	})
	getRect := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf(map[string]interface{}{"width": 800, "height": 600})
	})

	canvas.Set("getContext", getContext)
	canvas.Set("getBoundingClientRect", getRect)
	canvas.Set("width", 0)
	canvas.Set("height", 0)

	t.Cleanup(func() {
		for _, fn := range funcs {
			fn.Release()
		}
		getContext.Release()
		getRect.Release()
	})

	return canvas
}

func newZeroCanvas(t *testing.T) js.Value {
	canvas := js.ValueOf(map[string]interface{}{})
	ctx := js.ValueOf(map[string]interface{}{})
	setTransform := nopFunc()
	ctx.Set("setTransform", setTransform)
	getContext := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return ctx
	})
	getRect := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return js.ValueOf(map[string]interface{}{"width": 0, "height": 0})
	})

	canvas.Set("getContext", getContext)
	canvas.Set("getBoundingClientRect", getRect)
	canvas.Set("width", 0)
	canvas.Set("height", 0)

	t.Cleanup(func() {
		setTransform.Release()
		getContext.Release()
		getRect.Release()
	})

	return canvas
}

func nopFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return nil
	})
}
