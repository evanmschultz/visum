//go:build js && wasm

package web

import (
	"testing"
	"syscall/js"

	"github.com/evanschultz/visum/internal/app"
	"github.com/evanschultz/visum/internal/core"
)

func TestControllerSyncToDOM(t *testing.T) {
	js.Global().Set("document", js.ValueOf(map[string]interface{}{"activeElement": js.Null()}))

	engine := app.NewEngine(core.DefaultParams())
	controller := NewController(engine, nil)
	controller.elements = map[string]js.Value{
		"points":        newInput("0", false),
		"multiplier":    newInput("0", false),
		"rotation":      newInput("0", false),
		"start-index":   newInput("0", false),
		"line-count":    newInput("0", false),
		"line-count-all": newInput("", true),
		"show-circle":   newInput("", false),
		"show-points":   newInput("", false),
		"show-labels":   newInput("", false),
		"label-step":    newInput("0", false),
		"line-width":    newInput("0", false),
		"point-radius":  newInput("0", false),
		"bg-color":      newInput("#000000", false),
		"line-color":    newInput("#000000", false),
		"circle-color":  newInput("#000000", false),
		"point-color":   newInput("#000000", false),
		"label-color":   newInput("#000000", false),
		"step-amount":   newInput("0", false),
		"step-target":   newSelect("lines"),
		"line-anim-enable":  newInput("", false),
		"line-anim-start":   newInput("0", false),
		"line-anim-end":     newInput("0", false),
		"line-anim-speed":   newInput("0", false),
		"line-anim-loop":    newInput("", false),
		"line-anim-pingpong": newInput("", false),
		"mult-anim-enable":  newInput("", false),
		"mult-anim-start":   newInput("0", false),
		"mult-anim-end":     newInput("0", false),
		"mult-anim-speed":   newInput("0", false),
		"mult-anim-loop":    newInput("", false),
		"mult-anim-pingpong": newInput("", false),
		"points-anim-enable":  newInput("", false),
		"points-anim-start":   newInput("0", false),
		"points-anim-end":     newInput("0", false),
		"points-anim-speed":   newInput("0", false),
		"points-anim-loop":    newInput("", false),
		"points-anim-pingpong": newInput("", false),
		"play-toggle":   newInput("", false),
		"reverse-toggle": newInput("", false),
		"live-readout":  newInput("", false),
	}

	controller.SyncToDOM()

	if got := controller.elements["points"].Get("value").String(); got != "200" {
		t.Fatalf("expected points value 200, got %q", got)
	}
	if got := controller.elements["multiplier"].Get("value").String(); got != "2" {
		t.Fatalf("expected multiplier value 2, got %q", got)
	}
	if !controller.elements["line-count-all"].Get("checked").Bool() {
		t.Fatalf("expected line-count-all to be checked")
	}
	if got := controller.elements["play-toggle"].Get("textContent").String(); got != "PAUSE" {
		t.Fatalf("expected play button to read PAUSE, got %q", got)
	}
}

func TestControllerSyncFromDOM(t *testing.T) {
	js.Global().Set("document", js.ValueOf(map[string]interface{}{"activeElement": js.Null()}))

	engine := app.NewEngine(core.DefaultParams())
	controller := NewController(engine, nil)
	controller.elements = map[string]js.Value{
		"points":        newInput("120", false),
		"multiplier":    newInput("3.5", false),
		"rotation":      newInput("15", false),
		"start-index":   newInput("7", false),
		"line-count":    newInput("80", false),
		"line-count-all": newInput("", true),
		"show-circle":   newInput("", false),
		"show-points":   newInput("", true),
		"show-labels":   newInput("", true),
		"label-step":    newInput("5", false),
		"line-width":    newInput("2", false),
		"point-radius":  newInput("1", false),
		"bg-color":      newInput("#010101", false),
		"line-color":    newInput("#020202", false),
		"circle-color":  newInput("#030303", false),
		"point-color":   newInput("#040404", false),
		"label-color":   newInput("#050505", false),
		"step-amount":   newInput("2", false),
		"step-target":   newSelect("points"),
		"line-anim-enable":  newInput("", true),
		"line-anim-start":   newInput("0", false),
		"line-anim-end":     newInput("50", false),
		"line-anim-speed":   newInput("5", false),
		"line-anim-loop":    newInput("", true),
		"line-anim-pingpong": newInput("", false),
		"mult-anim-enable":  newInput("", true),
		"mult-anim-start":   newInput("1", false),
		"mult-anim-end":     newInput("4", false),
		"mult-anim-speed":   newInput("0.5", false),
		"mult-anim-loop":    newInput("", false),
		"mult-anim-pingpong": newInput("", true),
		"points-anim-enable":  newInput("", true),
		"points-anim-start":   newInput("10", false),
		"points-anim-end":     newInput("30", false),
		"points-anim-speed":   newInput("3", false),
		"points-anim-loop":    newInput("", false),
		"points-anim-pingpong": newInput("", true),
	}

	controller.SyncFromDOM()
	snapshot := engine.Snapshot()

	if snapshot.Params.PointCount != 120 {
		t.Fatalf("expected points 120, got %d", snapshot.Params.PointCount)
	}
	if snapshot.Params.Multiplier != 3.5 {
		t.Fatalf("expected multiplier 3.5, got %.2f", snapshot.Params.Multiplier)
	}
	if snapshot.Params.RotationDeg != 15 {
		t.Fatalf("expected rotation 15, got %.2f", snapshot.Params.RotationDeg)
	}
	if snapshot.Params.StartIndex != 7 {
		t.Fatalf("expected start index 7, got %d", snapshot.Params.StartIndex)
	}
	if snapshot.Step.Target != app.StepPoints {
		t.Fatalf("expected step target points")
	}
	if !snapshot.Animations.Lines.Settings.Enabled || snapshot.Animations.Lines.Settings.End != 50 {
		t.Fatalf("expected line animation settings applied")
	}
}

func TestControllerBinders(t *testing.T) {
	js.Global().Set("document", js.ValueOf(map[string]interface{}{"activeElement": js.Null()}))

	engine := app.NewEngine(core.DefaultParams())
	controller := NewController(engine, nil)

	handlerMap := map[string]js.Value{}
	numberEl := stubElement(t, "10", false, handlerMap)
	controller.elements = map[string]js.Value{"points": numberEl}

	controller.bindNumber("points", func(value float64) { engine.SetPointCount(int(value)) })
	numberEl.Set("value", "42")
	handlerMap["input"].Invoke()

	if engine.Snapshot().Params.PointCount != 42 {
		t.Fatalf("expected points 42, got %d", engine.Snapshot().Params.PointCount)
	}

	checkHandlers := map[string]js.Value{}
	checkEl := stubElement(t, "", true, checkHandlers)
	controller.elements["show-circle"] = checkEl
	controller.bindCheckbox("show-circle", func(checked bool) { engine.SetShowCircle(checked) })
	checkEl.Set("checked", false)
	checkHandlers["change"].Invoke()
	if engine.Snapshot().Params.ShowCircle {
		t.Fatalf("expected show circle false")
	}

	selectHandlers := map[string]js.Value{}
	selectEl := stubElement(t, "multiplier", false, selectHandlers)
	controller.elements["step-target"] = selectEl
	controller.bindSelect("step-target", func(value string) { controller.engine.SetStepTarget(app.StepMultiplier) })
	selectEl.Set("value", "multiplier")
	selectHandlers["change"].Invoke()
	if engine.Snapshot().Step.Target != app.StepMultiplier {
		t.Fatalf("expected step target multiplier")
	}

	buttonHandlers := map[string]js.Value{}
	buttonEl := stubElement(t, "", false, buttonHandlers)
	controller.elements["play-toggle"] = buttonEl
	controller.bindButton("play-toggle", func() { controller.engine.ToggleRunning() })
	engine.SetRunning(false)
	buttonHandlers["click"].Invoke()
	if !engine.Snapshot().Running {
		t.Fatalf("expected running true")
	}
}

func TestControllerBind(t *testing.T) {
	ids := []string{
		"points", "multiplier", "rotation", "start-index", "line-count", "line-count-all",
		"show-circle", "show-points", "show-labels", "label-step", "line-width", "point-radius",
		"bg-color", "line-color", "circle-color", "point-color", "label-color",
		"play-toggle", "reverse-toggle", "step-forward", "step-back", "step-target", "step-amount", "reset-params",
		"line-anim-enable", "line-anim-start", "line-anim-end", "line-anim-speed", "line-anim-loop", "line-anim-pingpong",
		"mult-anim-enable", "mult-anim-start", "mult-anim-end", "mult-anim-speed", "mult-anim-loop", "mult-anim-pingpong",
		"points-anim-enable", "points-anim-start", "points-anim-end", "points-anim-speed", "points-anim-loop", "points-anim-pingpong",
		"live-readout",
	}

	elements := make(map[string]js.Value, len(ids))
	for _, id := range ids {
		elements[id] = stubElementNoHandlers(t, "0", false)
	}
	setupDocument(t, elements)

	engine := app.NewEngine(core.DefaultParams())
	controller := NewController(engine, nil)
	controller.Bind()
}

func TestStepTargetValue(t *testing.T) {
	if stepTargetValue(app.StepMultiplier) != "multiplier" {
		t.Fatalf("expected multiplier step target")
	}
	if stepTargetValue(app.StepPoints) != "points" {
		t.Fatalf("expected points step target")
	}
	if stepTargetValue(app.StepLines) != "lines" {
		t.Fatalf("expected lines step target")
	}
}

func TestReadCheckboxAndFormatFloat(t *testing.T) {
	if !readCheckbox(newInput("", true)) {
		t.Fatalf("expected checkbox true")
	}
	if formatFloat(1.5) != "1.5" {
		t.Fatalf("unexpected float formatting")
	}
}

func TestReadoutFormatting(t *testing.T) {
	if got := formatNumber(12.3456, &CanvasRenderer{cssSize: core.Size{Width: 400}}); got != "12.35" {
		t.Fatalf("expected compact format, got %q", got)
	}
	if got := formatNumber(123.456, nil); got != "123.5" {
		t.Fatalf("expected reduced precision, got %q", got)
	}
	if got := joinParts([]string{"k=2", "N=200"}); got != "k=2 | N=200" {
		t.Fatalf("unexpected join output: %q", got)
	}
}

func TestUpdateReadout(t *testing.T) {
	engine := app.NewEngine(core.DefaultParams())
	engine.SetMultiplierAnimation(app.AnimationSettings{Enabled: true})
	engine.SetLineAnimation(app.AnimationSettings{Enabled: true})

	controller := NewController(engine, &CanvasRenderer{cssSize: core.Size{Width: 900}})
	controller.elements = map[string]js.Value{
		"live-readout": newInput("", false),
	}

	controller.updateReadout(engine.Snapshot())
	if got := controller.elements["live-readout"].Get("textContent").String(); got == "" {
		t.Fatalf("expected readout to be populated")
	}

	engine.SetMultiplierAnimation(app.AnimationSettings{Enabled: false})
	engine.SetLineAnimation(app.AnimationSettings{Enabled: false})
	engine.SetPointAnimation(app.AnimationSettings{Enabled: false})
	controller.updateReadout(engine.Snapshot())
	if got := controller.elements["live-readout"].Get("textContent").String(); got == "" {
		t.Fatalf("expected readout to include k")
	}
}

func TestControllerSetters(t *testing.T) {
	js.Global().Set("document", js.ValueOf(map[string]interface{}{"activeElement": js.Null()}))

	engine := app.NewEngine(core.DefaultParams())
	controller := NewController(engine, nil)
	controller.elements = map[string]js.Value{
		"points":        newInput("0", false),
		"line-count-all": newInput("", false),
		"bg-color":      newInput("#000000", false),
		"step-target":   newSelect("lines"),
		"live-readout":  newInput("", false),
	}

	controller.setInputValue("points", 12)
	controller.setCheckbox("line-count-all", true)
	controller.setColorValue("bg-color", "#ffffff")
	controller.setSelectValue("step-target", "multiplier")

	if controller.elements["points"].Get("value").String() != "12" {
		t.Fatalf("expected points to be set")
	}
	if !controller.elements["line-count-all"].Get("checked").Bool() {
		t.Fatalf("expected checkbox to be set")
	}
	if controller.elements["bg-color"].Get("value").String() != "#ffffff" {
		t.Fatalf("expected color value updated")
	}
	if controller.elements["step-target"].Get("value").String() != "multiplier" {
		t.Fatalf("expected select value updated")
	}
}

func TestBindColor(t *testing.T) {
	js.Global().Set("document", js.ValueOf(map[string]interface{}{"activeElement": js.Null()}))
	engine := app.NewEngine(core.DefaultParams())
	controller := NewController(engine, nil)

	handlers := map[string]js.Value{}
	colorEl := stubElement(t, "#123456", false, handlers)
	controller.elements = map[string]js.Value{"bg-color": colorEl}

	controller.bindColor("bg-color", func(value string) { engine.SetBackgroundColor(value) })
	colorEl.Set("value", "#abcdef")
	handlers["input"].Invoke()

	if engine.Snapshot().Params.Colors.Background != "#abcdef" {
		t.Fatalf("expected background color to update")
	}
}

func TestBindStepButtonHold(t *testing.T) {
	js.Global().Set("document", js.ValueOf(map[string]interface{}{"activeElement": js.Null()}))

	engine := app.NewEngine(core.DefaultParams())
	engine.SetLineCount(10)
	controller := NewController(engine, nil)

	handlers := map[string]js.Value{}
	stepEl := stubElement(t, "", false, handlers)
	controller.elements = map[string]js.Value{"step-forward": stepEl}

	var timeoutFn js.Value
	var intervalFn js.Value

	setTimeout := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		timeoutFn = args[0]
		return js.ValueOf(1)
	})
	setInterval := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		intervalFn = args[0]
		return js.ValueOf(2)
	})
	clearTimeout := js.FuncOf(func(this js.Value, args []js.Value) interface{} { return nil })
	clearInterval := js.FuncOf(func(this js.Value, args []js.Value) interface{} { return nil })

	js.Global().Set("setTimeout", setTimeout)
	js.Global().Set("setInterval", setInterval)
	js.Global().Set("clearTimeout", clearTimeout)
	js.Global().Set("clearInterval", clearInterval)

	t.Cleanup(func() {
		setTimeout.Release()
		setInterval.Release()
		clearTimeout.Release()
		clearInterval.Release()
	})

	controller.bindStepButton("step-forward", 1)

	handlers["click"].Invoke()
	if engine.Snapshot().Params.LineCount != 11 {
		t.Fatalf("expected click to step once")
	}

	prevent := js.FuncOf(func(this js.Value, args []js.Value) interface{} { return nil })
	t.Cleanup(func() { prevent.Release() })
	event := js.ValueOf(map[string]interface{}{"preventDefault": prevent})
	handlers["pointerdown"].Invoke(event)
	if !timeoutFn.Truthy() {
		t.Fatalf("expected hold timeout to be set")
	}
	timeoutFn.Invoke()
	if !intervalFn.Truthy() {
		t.Fatalf("expected interval to be set")
	}
	intervalFn.Invoke()
	handlers["pointerup"].Invoke(event)

	if engine.Snapshot().Params.LineCount < 13 {
		t.Fatalf("expected hold to step multiple times")
	}
}

func TestSyncAnimationMissingElements(t *testing.T) {
	controller := NewController(app.NewEngine(core.DefaultParams()), nil)
	controller.elements = map[string]js.Value{}
	controller.syncAnimation("line-anim", func(settings app.AnimationSettings) {
		t.Fatalf("expected syncAnimation to return early")
	})
}

func TestAnimationBindings(t *testing.T) {
	js.Global().Set("document", js.ValueOf(map[string]interface{}{"activeElement": js.Null()}))

	engine := app.NewEngine(core.DefaultParams())
	controller := NewController(engine, nil)

	handlers := map[string]map[string]js.Value{}
	controller.elements = map[string]js.Value{
		"line-anim-enable":  stubElement(t, "", true, newHandlerMap(handlers, "enable")),
		"line-anim-start":   stubElement(t, "0", false, newHandlerMap(handlers, "start")),
		"line-anim-end":     stubElement(t, "10", false, newHandlerMap(handlers, "end")),
		"line-anim-speed":   stubElement(t, "2", false, newHandlerMap(handlers, "speed")),
		"line-anim-loop":    stubElement(t, "", true, newHandlerMap(handlers, "loop")),
		"line-anim-pingpong": stubElement(t, "", false, newHandlerMap(handlers, "pingpong")),
	}

	var applied app.AnimationSettings
	controller.bindAnimationInputs("line-anim", func(settings app.AnimationSettings) { applied = settings })

	handlers["enable"]["input"].Invoke()

	if !applied.Enabled || applied.End != 10 || !applied.Loop {
		t.Fatalf("expected animation settings applied, got %+v", applied)
	}
}

func TestSetInputValueActiveElement(t *testing.T) {
	el := newInput("1", false)
	doc := js.ValueOf(map[string]interface{}{"activeElement": el})
	js.Global().Set("document", doc)

	controller := NewController(app.NewEngine(core.DefaultParams()), nil)
	controller.elements = map[string]js.Value{"points": el}
	controller.setInputValue("points", 99)

	if el.Get("value").String() != "1" {
		t.Fatalf("expected active element value unchanged")
	}
}

func TestReadFloatInvalid(t *testing.T) {
	el := newInput("nope", false)
	if readFloat(el) != 0 {
		t.Fatalf("expected invalid float to return 0")
	}
}

func newInput(value string, checked bool) js.Value {
	el := js.ValueOf(map[string]interface{}{})
	el.Set("value", value)
	el.Set("checked", checked)
	el.Set("textContent", "")
	return el
}

func newSelect(value string) js.Value {
	el := js.ValueOf(map[string]interface{}{})
	el.Set("value", value)
	return el
}

func stubElement(t *testing.T, value string, checked bool, handlers map[string]js.Value) js.Value {
	el := newInput(value, checked)
	add := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) >= 2 {
			handlers[args[0].String()] = args[1]
		}
		return nil
	})
	el.Set("addEventListener", add)
	t.Cleanup(func() { add.Release() })
	return el
}

func stubElementNoHandlers(t *testing.T, value string, checked bool) js.Value {
	el := newInput(value, checked)
	add := js.FuncOf(func(this js.Value, args []js.Value) interface{} { return nil })
	el.Set("addEventListener", add)
	t.Cleanup(func() { add.Release() })
	return el
}

func newHandlerMap(container map[string]map[string]js.Value, key string) map[string]js.Value {
	if container[key] == nil {
		container[key] = map[string]js.Value{}
	}
	return container[key]
}
