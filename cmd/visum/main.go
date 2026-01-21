//go:build js && wasm

package main

import (
	"runtime"

	"github.com/evanschultz/visum/internal/adapter/web"
	"github.com/evanschultz/visum/internal/app"
	"github.com/evanschultz/visum/internal/core"
)

func main() {
	runtime.LockOSThread()

	engine := app.NewEngine(core.DefaultParams())
	renderer, err := web.NewCanvasRenderer("visum-canvas")
	if err != nil {
		return
	}

	controller := web.NewController(engine, renderer)
	controller.Bind()
	web.StartLoop(engine, renderer, controller)

	select {}
}
