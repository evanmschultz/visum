//go:build js && wasm

package web

import (
	"math"
	"strconv"
	"syscall/js"

	"github.com/evanschultz/visum/internal/app"
	"github.com/evanschultz/visum/internal/core"
)

// Controller wires DOM inputs to the engine state.
type Controller struct {
	doc        js.Value
	engine     *app.Engine
	renderer   *CanvasRenderer
	elements   map[string]js.Value
	callbacks  []js.Func
	holdStates map[string]*holdState
	reverse    bool
}

type holdState struct {
	holdTimeout  js.Value
	interval     js.Value
	holding      bool
	consumeClick bool
}

// NewController creates a controller for the UI.
func NewController(engine *app.Engine, renderer *CanvasRenderer) *Controller {
	return &Controller{
		doc:        js.Global().Get("document"),
		engine:     engine,
		renderer:   renderer,
		elements:   make(map[string]js.Value),
		holdStates: make(map[string]*holdState),
	}
}

// Bind registers DOM event handlers and syncs initial state.
func (c *Controller) Bind() {
	c.cacheElements([]string{
		"points", "multiplier", "rotation", "start-index", "line-count", "line-count-all",
		"show-circle", "show-points", "show-labels", "label-step", "line-width", "point-radius",
		"bg-color", "line-color", "circle-color", "point-color", "label-color",
		"play-toggle", "reverse-toggle", "step-forward", "step-back", "step-target", "step-amount", "reset-params",
		"line-anim-enable", "line-anim-start", "line-anim-end", "line-anim-speed", "line-anim-loop", "line-anim-pingpong",
		"mult-anim-enable", "mult-anim-start", "mult-anim-end", "mult-anim-speed", "mult-anim-loop", "mult-anim-pingpong",
		"points-anim-enable", "points-anim-start", "points-anim-end", "points-anim-speed", "points-anim-loop", "points-anim-pingpong",
		"live-readout",
	})

	c.bindNumber("points", func(value float64) { c.engine.SetPointCount(int(value)) })
	c.bindNumber("multiplier", func(value float64) { c.engine.SetMultiplier(value) })
	c.bindNumber("rotation", func(value float64) { c.engine.SetRotationDeg(value) })
	c.bindNumber("start-index", func(value float64) { c.engine.SetStartIndex(int(value)) })
	c.bindNumber("line-count", func(value float64) { c.engine.SetLineCount(int(value)) })
	c.bindCheckbox("line-count-all", func(checked bool) { c.engine.SetLineAll(checked) })

	c.bindCheckbox("show-circle", func(checked bool) { c.engine.SetShowCircle(checked) })
	c.bindCheckbox("show-points", func(checked bool) { c.engine.SetShowPoints(checked) })
	c.bindCheckbox("show-labels", func(checked bool) { c.engine.SetShowLabels(checked) })
	c.bindNumber("label-step", func(value float64) { c.engine.SetLabelStep(int(value)) })
	c.bindNumber("line-width", func(value float64) { c.engine.SetLineWidth(value) })
	c.bindNumber("point-radius", func(value float64) { c.engine.SetPointRadius(value) })

	c.bindColor("bg-color", func(value string) { c.engine.SetBackgroundColor(value) })
	c.bindColor("line-color", func(value string) { c.engine.SetLineColor(value) })
	c.bindColor("circle-color", func(value string) { c.engine.SetCircleColor(value) })
	c.bindColor("point-color", func(value string) { c.engine.SetPointColor(value) })
	c.bindColor("label-color", func(value string) { c.engine.SetLabelColor(value) })

	c.bindButton("play-toggle", func() { c.engine.ToggleRunning() })
	c.bindButton("reverse-toggle", func() {
		c.reverse = !c.reverse
		c.engine.SetReverse(c.reverse)
	})
	c.bindStepButton("step-forward", 1)
	c.bindStepButton("step-back", -1)
	c.bindButton("reset-params", func() { c.resetToDefaults() })

	c.bindSelect("step-target", func(value string) {
		switch value {
		case "multiplier":
			c.engine.SetStepTarget(app.StepMultiplier)
		case "points":
			c.engine.SetStepTarget(app.StepPoints)
		default:
			c.engine.SetStepTarget(app.StepLines)
		}
	})
	c.bindNumber("step-amount", func(value float64) { c.engine.SetStepAmount(value) })

	c.bindAnimation("line-anim", func(settings app.AnimationSettings) { c.engine.SetLineAnimation(settings) })
	c.bindAnimation("mult-anim", func(settings app.AnimationSettings) { c.engine.SetMultiplierAnimation(settings) })
	c.bindAnimation("points-anim", func(settings app.AnimationSettings) { c.engine.SetPointAnimation(settings) })

	c.SyncFromDOM()
	c.engine.SetRunning(true)
	c.SyncToDOM()
}

// SyncFromDOM pulls the current UI values into the engine.
func (c *Controller) SyncFromDOM() {
	c.syncNumber("points", func(v float64) { c.engine.SetPointCount(int(v)) })
	c.syncNumber("multiplier", func(v float64) { c.engine.SetMultiplier(v) })
	c.syncNumber("rotation", func(v float64) { c.engine.SetRotationDeg(v) })
	c.syncNumber("start-index", func(v float64) { c.engine.SetStartIndex(int(v)) })
	c.syncNumber("line-count", func(v float64) { c.engine.SetLineCount(int(v)) })
	c.syncCheckbox("line-count-all", func(v bool) { c.engine.SetLineAll(v) })
	c.syncCheckbox("show-circle", func(v bool) { c.engine.SetShowCircle(v) })
	c.syncCheckbox("show-points", func(v bool) { c.engine.SetShowPoints(v) })
	c.syncCheckbox("show-labels", func(v bool) { c.engine.SetShowLabels(v) })
	c.syncNumber("label-step", func(v float64) { c.engine.SetLabelStep(int(v)) })
	c.syncNumber("line-width", func(v float64) { c.engine.SetLineWidth(v) })
	c.syncNumber("point-radius", func(v float64) { c.engine.SetPointRadius(v) })
	c.syncColor("bg-color", func(v string) { c.engine.SetBackgroundColor(v) })
	c.syncColor("line-color", func(v string) { c.engine.SetLineColor(v) })
	c.syncColor("circle-color", func(v string) { c.engine.SetCircleColor(v) })
	c.syncColor("point-color", func(v string) { c.engine.SetPointColor(v) })
	c.syncColor("label-color", func(v string) { c.engine.SetLabelColor(v) })

	c.syncNumber("step-amount", func(v float64) { c.engine.SetStepAmount(v) })
	c.syncSelect("step-target", func(v string) {
		switch v {
		case "multiplier":
			c.engine.SetStepTarget(app.StepMultiplier)
		case "points":
			c.engine.SetStepTarget(app.StepPoints)
		default:
			c.engine.SetStepTarget(app.StepLines)
		}
	})
	c.syncAnimation("line-anim", func(settings app.AnimationSettings) { c.engine.SetLineAnimation(settings) })
	c.syncAnimation("mult-anim", func(settings app.AnimationSettings) { c.engine.SetMultiplierAnimation(settings) })
	c.syncAnimation("points-anim", func(settings app.AnimationSettings) { c.engine.SetPointAnimation(settings) })
}

// SyncToDOM updates UI values from the engine state for live animations.
func (c *Controller) SyncToDOM() {
	snapshot := c.engine.Snapshot()
	params := snapshot.Params

	c.setInputValue("points", float64(params.PointCount))
	c.setInputValue("multiplier", params.Multiplier)
	c.setInputValue("rotation", params.RotationDeg)
	c.setInputValue("start-index", float64(params.StartIndex))
	if params.LineCount < 0 {
		c.setCheckbox("line-count-all", true)
		c.setInputValue("line-count", float64(params.PointCount))
	} else {
		c.setCheckbox("line-count-all", false)
		c.setInputValue("line-count", float64(params.LineCount))
	}

	c.setCheckbox("show-circle", params.ShowCircle)
	c.setCheckbox("show-points", params.ShowPoints)
	c.setCheckbox("show-labels", params.ShowLabels)
	c.setInputValue("label-step", float64(params.LabelStep))
	c.setInputValue("line-width", params.LineWidth)
	c.setInputValue("point-radius", params.PointRadius)

	c.setColorValue("bg-color", params.Colors.Background)
	c.setColorValue("line-color", params.Colors.Line)
	c.setColorValue("circle-color", params.Colors.Circle)
	c.setColorValue("point-color", params.Colors.Point)
	c.setColorValue("label-color", params.Colors.Label)

	c.setInputValue("step-amount", snapshot.Step.Amount)
	c.setSelectValue("step-target", stepTargetValue(snapshot.Step.Target))

	c.setAnimationInputs("line-anim", snapshot.Animations.Lines.Settings)
	c.setAnimationInputs("mult-anim", snapshot.Animations.Multiplier.Settings)
	c.setAnimationInputs("points-anim", snapshot.Animations.Points.Settings)

	playLabel := "PLAY"
	if snapshot.Running {
		playLabel = "PAUSE"
	}
	if el, ok := c.elements["play-toggle"]; ok {
		el.Set("textContent", playLabel)
	}

	reverseLabel := "REVERSE"
	if c.reverse {
		reverseLabel = "FORWARD"
	}
	if el, ok := c.elements["reverse-toggle"]; ok {
		el.Set("textContent", reverseLabel)
	}

	c.updateReadout(snapshot)
}

func (c *Controller) cacheElements(ids []string) {
	for _, id := range ids {
		el := c.doc.Call("getElementById", id)
		if !el.IsNull() && !el.IsUndefined() {
			c.elements[id] = el
		}
	}
}

func (c *Controller) bindNumber(id string, apply func(value float64)) {
	el, ok := c.elements[id]
	if !ok {
		return
	}
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		apply(readFloat(el))
		return nil
	})
	el.Call("addEventListener", "input", cb)
	c.callbacks = append(c.callbacks, cb)
}

func (c *Controller) bindCheckbox(id string, apply func(checked bool)) {
	el, ok := c.elements[id]
	if !ok {
		return
	}
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		apply(el.Get("checked").Bool())
		return nil
	})
	el.Call("addEventListener", "change", cb)
	c.callbacks = append(c.callbacks, cb)
}

func (c *Controller) bindColor(id string, apply func(value string)) {
	el, ok := c.elements[id]
	if !ok {
		return
	}
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		apply(el.Get("value").String())
		return nil
	})
	el.Call("addEventListener", "input", cb)
	c.callbacks = append(c.callbacks, cb)
}

func (c *Controller) bindButton(id string, apply func()) {
	el, ok := c.elements[id]
	if !ok {
		return
	}
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		apply()
		return nil
	})
	el.Call("addEventListener", "click", cb)
	c.callbacks = append(c.callbacks, cb)
}

func (c *Controller) bindStepButton(id string, direction int) {
	el, ok := c.elements[id]
	if !ok {
		return
	}

	state := &holdState{}
	c.holdStates[id] = state

	startRepeat := func() {
		state.holding = true
		state.consumeClick = true
		c.engine.Step(direction)
		intervalFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			c.engine.Step(direction)
			return nil
		})
		state.interval = js.Global().Call("setInterval", intervalFunc, 60)
		c.callbacks = append(c.callbacks, intervalFunc)
	}

	start := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 && args[0].Truthy() && args[0].Get("pointerId").Truthy() {
			el.Call("setPointerCapture", args[0].Get("pointerId"))
		}
		if state.holdTimeout.Truthy() || state.interval.Truthy() {
			return nil
		}
		state.consumeClick = false
		timeoutFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			startRepeat()
			return nil
		})
		state.holdTimeout = js.Global().Call("setTimeout", timeoutFunc, 150)
		c.callbacks = append(c.callbacks, timeoutFunc)
		return nil
	})

	stop := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if state.holdTimeout.Truthy() {
			js.Global().Call("clearTimeout", state.holdTimeout)
			state.holdTimeout = js.Undefined()
		}
		if state.interval.Truthy() {
			js.Global().Call("clearInterval", state.interval)
			state.interval = js.Undefined()
		}
		state.holding = false
		return nil
	})

	click := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if state.consumeClick {
			state.consumeClick = false
			return nil
		}
		c.engine.Step(direction)
		return nil
	})

	el.Call("addEventListener", "pointerdown", start)
	el.Call("addEventListener", "pointerup", stop)
	el.Call("addEventListener", "pointerleave", stop)
	el.Call("addEventListener", "pointercancel", stop)
	el.Call("addEventListener", "lostpointercapture", stop)
	el.Call("addEventListener", "mousedown", start)
	el.Call("addEventListener", "mouseup", stop)
	el.Call("addEventListener", "mouseleave", stop)
	el.Call("addEventListener", "touchstart", start)
	el.Call("addEventListener", "touchend", stop)
	el.Call("addEventListener", "touchcancel", stop)
	el.Call("addEventListener", "click", click)

	c.callbacks = append(c.callbacks, start, stop, click)
}

func (c *Controller) resetToDefaults() {
	c.resetInputsToDefault()
	c.engine.Reset(core.DefaultParams())
	c.reverse = false
	c.engine.SetReverse(false)
	c.SyncFromDOM()
	c.SyncToDOM()
}

func (c *Controller) resetInputsToDefault() {
	for _, el := range c.elements {
		if el.IsNull() || el.IsUndefined() {
			continue
		}
		switch el.Get("tagName").String() {
		case "INPUT":
			if el.Get("type").String() == "checkbox" {
				el.Set("checked", el.Get("defaultChecked"))
			} else {
				el.Set("value", el.Get("defaultValue"))
			}
		case "SELECT":
			el.Set("value", el.Get("defaultValue"))
		}
	}
}

func (c *Controller) bindSelect(id string, apply func(value string)) {
	el, ok := c.elements[id]
	if !ok {
		return
	}
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		apply(el.Get("value").String())
		return nil
	})
	el.Call("addEventListener", "change", cb)
	c.callbacks = append(c.callbacks, cb)
}

func (c *Controller) bindAnimation(prefix string, apply func(settings app.AnimationSettings)) {
	c.bindAnimationInputs(prefix, apply)
}

func (c *Controller) bindAnimationInputs(prefix string, apply func(settings app.AnimationSettings)) {
	ids := map[string]string{
		"enable":   prefix + "-enable",
		"start":    prefix + "-start",
		"end":      prefix + "-end",
		"speed":    prefix + "-speed",
		"loop":     prefix + "-loop",
		"pingpong": prefix + "-pingpong",
	}

	applySettings := func() {
		settings := app.AnimationSettings{}
		settings.Enabled = readCheckbox(c.elements[ids["enable"]])
		settings.Start = readFloat(c.elements[ids["start"]])
		settings.End = readFloat(c.elements[ids["end"]])
		settings.Speed = readFloat(c.elements[ids["speed"]])
		settings.Loop = readCheckbox(c.elements[ids["loop"]])
		settings.PingPong = readCheckbox(c.elements[ids["pingpong"]])
		apply(settings)
	}

	for _, id := range ids {
		el, ok := c.elements[id]
		if !ok {
			continue
		}
		cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			applySettings()
			return nil
		})
		el.Call("addEventListener", "input", cb)
		el.Call("addEventListener", "change", cb)
		c.callbacks = append(c.callbacks, cb)
	}
}

func (c *Controller) syncNumber(id string, apply func(value float64)) {
	if el, ok := c.elements[id]; ok {
		apply(readFloat(el))
	}
}

func (c *Controller) syncCheckbox(id string, apply func(value bool)) {
	if el, ok := c.elements[id]; ok {
		apply(el.Get("checked").Bool())
	}
}

func (c *Controller) syncColor(id string, apply func(value string)) {
	if el, ok := c.elements[id]; ok {
		apply(el.Get("value").String())
	}
}

func (c *Controller) syncSelect(id string, apply func(value string)) {
	if el, ok := c.elements[id]; ok {
		apply(el.Get("value").String())
	}
}

func (c *Controller) syncAnimation(prefix string, apply func(settings app.AnimationSettings)) {
	ids := []string{prefix + "-enable", prefix + "-start", prefix + "-end", prefix + "-speed", prefix + "-loop", prefix + "-pingpong"}
	for _, id := range ids {
		if _, ok := c.elements[id]; !ok {
			return
		}
	}
	settings := app.AnimationSettings{
		Enabled:  readCheckbox(c.elements[prefix+"-enable"]),
		Start:    readFloat(c.elements[prefix+"-start"]),
		End:      readFloat(c.elements[prefix+"-end"]),
		Speed:    readFloat(c.elements[prefix+"-speed"]),
		Loop:     readCheckbox(c.elements[prefix+"-loop"]),
		PingPong: readCheckbox(c.elements[prefix+"-pingpong"]),
	}
	apply(settings)
}

func (c *Controller) setInputValue(id string, value float64) {
	el, ok := c.elements[id]
	if !ok {
		return
	}
	if isActiveElement(el) {
		return
	}
	el.Set("value", formatFloat(value))
}

func (c *Controller) setCheckbox(id string, value bool) {
	el, ok := c.elements[id]
	if !ok {
		return
	}
	if isActiveElement(el) {
		return
	}
	el.Set("checked", value)
}

func (c *Controller) setColorValue(id string, value string) {
	el, ok := c.elements[id]
	if !ok {
		return
	}
	el.Set("value", value)
}

func (c *Controller) setSelectValue(id string, value string) {
	el, ok := c.elements[id]
	if !ok {
		return
	}
	if isActiveElement(el) {
		return
	}
	el.Set("value", value)
}

func (c *Controller) setAnimationInputs(prefix string, settings app.AnimationSettings) {
	c.setCheckbox(prefix+"-enable", settings.Enabled)
	c.setInputValue(prefix+"-start", settings.Start)
	c.setInputValue(prefix+"-end", settings.End)
	c.setInputValue(prefix+"-speed", settings.Speed)
	c.setCheckbox(prefix+"-loop", settings.Loop)
	c.setCheckbox(prefix+"-pingpong", settings.PingPong)
}

func (c *Controller) updateReadout(snapshot app.Snapshot) {
	el, ok := c.elements["live-readout"]
	if !ok {
		return
	}
	parts := make([]string, 0, 3)
	parts = append(parts, "k="+formatNumber(snapshot.Params.Multiplier, c.renderer))
	if snapshot.Animations.Points.Settings.Enabled {
		parts = append(parts, "N="+formatInt(snapshot.Params.PointCount))
	}
	if snapshot.Animations.Lines.Settings.Enabled {
		lines := snapshot.Params.LineCount
		if lines < 0 {
			lines = snapshot.Params.PointCount
		}
		parts = append(parts, "LINES="+formatInt(lines))
	}
	el.Set("textContent", joinParts(parts))
}

func readFloat(el js.Value) float64 {
	value := el.Get("value").String()
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func readCheckbox(el js.Value) bool {
	if el.IsNull() || el.IsUndefined() {
		return false
	}
	return el.Get("checked").Bool()
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatInt(value int) string {
	return strconv.Itoa(value)
}

func formatNumber(value float64, renderer *CanvasRenderer) string {
	width := 800.0
	if renderer != nil {
		if size := renderer.Size(); size.Width > 0 {
			width = size.Width
		}
	}
	precision := 3
	if width < 520 {
		precision = 2
	}
	if math.Abs(value) >= 100 {
		precision = 1
	}
	return strconv.FormatFloat(value, 'f', precision, 64)
}

func joinParts(parts []string) string {
	if len(parts) == 1 {
		return parts[0]
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += " | " + parts[i]
	}
	return result
}

func stepTargetValue(target app.StepTarget) string {
	switch target {
	case app.StepMultiplier:
		return "multiplier"
	case app.StepPoints:
		return "points"
	default:
		return "lines"
	}
}

func isActiveElement(el js.Value) bool {
	doc := js.Global().Get("document")
	if doc.IsUndefined() || doc.IsNull() {
		return false
	}
	active := doc.Get("activeElement")
	return active.Equal(el)
}
