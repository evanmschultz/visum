//go:build js && wasm

package web

import (
	"syscall/js"

	"github.com/evanschultz/visum/internal/app"
)

// StartLoop begins the requestAnimationFrame render loop.
func StartLoop(engine *app.Engine, renderer *CanvasRenderer, controller *Controller) {
	var last float64
	var raf js.Func

	raf = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		now := args[0].Float()
		if last == 0 {
			last = now
		}
		dt := (now - last) / 1000
		last = now

		engine.Update(dt)
		frame := engine.Frame(renderer.Size())
		renderer.Render(frame, engine.Snapshot().Params)
		controller.SyncToDOM()

		js.Global().Call("requestAnimationFrame", raf)
		return nil
	})

	js.Global().Call("requestAnimationFrame", raf)
}
