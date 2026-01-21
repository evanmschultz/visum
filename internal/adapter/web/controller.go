//go:build js && wasm

package web

import (
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
}

// NewController creates a controller for the UI.
func NewController(engine *app.Engine, renderer *CanvasRenderer) *Controller {
	return &Controller{
		doc:      js.Global().Get("document"),
		engine:   engine,
		renderer: renderer,
		elements: make(map[string]js.Value),
	}
}

// Bind registers DOM event handlers and syncs initial state.
func (c *Controller) Bind() {
	c.cacheElements([]string{
		"points", "multiplier", "rotation", "start-index", "line-count", "line-count-all",
		"show-circle", "show-points", "show-labels", "label-step", "line-width", "point-radius",
		"bg-color", "line-color", "circle-color", "point-color", "label-color",
		"play-toggle", "step-forward", "step-back", "step-target", "step-amount", "reset-params",
		"line-anim-enable", "line-anim-start", "line-anim-end", "line-anim-speed", "line-anim-loop", "line-anim-pingpong",
		"mult-anim-enable", "mult-anim-start", "mult-anim-end", "mult-anim-speed", "mult-anim-loop", "mult-anim-pingpong",
		"points-anim-enable", "points-anim-start", "points-anim-end", "points-anim-speed", "points-anim-loop", "points-anim-pingpong",
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
	c.bindButton("step-forward", func() { c.engine.Step(1) })
	c.bindButton("step-back", func() { c.engine.Step(-1) })
	c.bindButton("reset-params", func() { c.engine.Reset(core.DefaultParams()) })

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

	playLabel := "Play"
	if snapshot.Running {
		playLabel = "Pause"
	}
	if el, ok := c.elements["play-toggle"]; ok {
		el.Set("textContent", playLabel)
	}
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
