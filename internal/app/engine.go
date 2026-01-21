package app

import (
	"math"

	"github.com/evanschultz/visum/internal/core"
)

// StepTarget controls which parameter is advanced by a step action.
type StepTarget int

const (
	StepLines StepTarget = iota
	StepMultiplier
	StepPoints
)

// StepConfig controls how manual stepping behaves.
type StepConfig struct {
	Target StepTarget
	Amount float64
}

// AnimationSettings define a user-configurable animation track.
type AnimationSettings struct {
	Enabled  bool
	Start    float64
	End      float64
	Speed    float64
	Loop     bool
	PingPong bool
}

// Animation tracks the live animation state for a parameter.
type Animation struct {
	Settings AnimationSettings
	Value    float64
	Forward  bool
}

// Animations groups all animated parameters.
type Animations struct {
	Lines      Animation
	Multiplier Animation
	Points     Animation
}

// Snapshot captures the engine state for UI sync.
type Snapshot struct {
	Params     core.Params
	Animations Animations
	Running    bool
	Step       StepConfig
}

// Engine owns the current state, animations, and frame generation.
type Engine struct {
	params     core.Params
	animations Animations
	running    bool
	step       StepConfig
	reverse    bool
}

// NewEngine creates a new engine with default settings.
func NewEngine(params core.Params) *Engine {
	engine := &Engine{
		params: core.NormalizeParams(params),
		step: StepConfig{
			Target: StepLines,
			Amount: 1,
		},
		running: true,
	}
	engine.animations = Animations{
		Lines: Animation{Settings: AnimationSettings{Enabled: false, Start: 0, End: float64(engine.params.PointCount), Speed: 60}},
		Multiplier: Animation{Settings: AnimationSettings{Enabled: false, Start: engine.params.Multiplier, End: engine.params.Multiplier + 5, Speed: 0.2}},
		Points: Animation{Settings: AnimationSettings{Enabled: false, Start: float64(engine.params.PointCount), End: float64(engine.params.PointCount), Speed: 1}},
	}
	engine.animations.Lines.Value = float64(engine.params.PointCount)
	engine.animations.Multiplier.Value = engine.params.Multiplier
	engine.animations.Points.Value = float64(engine.params.PointCount)
	engine.animations.Lines.Forward = true
	engine.animations.Multiplier.Forward = true
	engine.animations.Points.Forward = true
	return engine
}

// Reset replaces the current parameters with the provided defaults.
func (e *Engine) Reset(params core.Params) {
	replacement := NewEngine(params)
	*e = *replacement
}

// Snapshot returns a copy of the current engine state.
func (e *Engine) Snapshot() Snapshot {
	return Snapshot{
		Params:     e.params,
		Animations: e.animations,
		Running:    e.running,
		Step:       e.step,
	}
}

// SetRunning sets the animation running state.
func (e *Engine) SetRunning(running bool) {
	e.running = running
}

// ToggleRunning flips the animation running state.
func (e *Engine) ToggleRunning() {
	e.running = !e.running
}

// SetReverse sets the animation direction (false = forward, true = reverse).
func (e *Engine) SetReverse(reverse bool) {
	e.reverse = reverse
}

// SetStepTarget sets the target for manual stepping.
func (e *Engine) SetStepTarget(target StepTarget) {
	e.step.Target = target
}

// SetStepAmount sets the amount for manual stepping.
func (e *Engine) SetStepAmount(amount float64) {
	if amount == 0 {
		amount = 1
	}
	e.step.Amount = amount
}

// Step advances or rewinds a parameter by the configured step amount.
func (e *Engine) Step(direction int) {
	if direction == 0 {
		return
	}
	amount := e.step.Amount
	if amount == 0 {
		amount = 1
	}

	switch e.step.Target {
	case StepMultiplier:
		e.SetMultiplier(e.params.Multiplier + float64(direction)*amount)
	case StepPoints:
		step := int(math.Round(amount))
		if step == 0 {
			step = 1
		}
		e.SetPointCount(e.params.PointCount + direction*step)
	case StepLines:
		step := int(math.Round(amount))
		if step == 0 {
			step = 1
		}
		if e.params.LineCount < 0 {
			e.params.LineCount = e.params.PointCount
		}
		e.SetLineCount(e.params.LineCount + direction*step)
	}
}

// Update advances animations based on the elapsed time in seconds.
func (e *Engine) Update(dt float64) {
	if !e.running {
		return
	}
	if dt == 0 {
		return
	}

	if e.reverse {
		dt = -dt
	}

	if e.animations.Lines.Settings.Enabled {
		value := e.animations.Lines.Advance(dt)
		e.SetLineCount(int(math.Round(value)))
	}
	if e.animations.Multiplier.Settings.Enabled {
		value := e.animations.Multiplier.Advance(dt)
		e.SetMultiplier(value)
	}
	if e.animations.Points.Settings.Enabled {
		value := e.animations.Points.Advance(dt)
		e.SetPointCount(int(math.Round(value)))
	}
}

// Frame returns the current geometry for rendering.
func (e *Engine) Frame(size core.Size) core.Frame {
	return core.BuildFrame(e.params, size)
}

// SetPointCount updates the number of points on the circle.
func (e *Engine) SetPointCount(count int) {
	if count < 2 {
		count = 2
	}
	if count > 4000 {
		count = 4000
	}
	e.params.PointCount = count
	if e.params.LineCount > count {
		e.params.LineCount = count
	}
}

// SetMultiplier updates the multiplier.
func (e *Engine) SetMultiplier(multiplier float64) {
	e.params.Multiplier = multiplier
}

// SetRotationDeg updates the rotation in degrees.
func (e *Engine) SetRotationDeg(deg float64) {
	e.params.RotationDeg = deg
}

// SetStartIndex updates the starting index for line drawing.
func (e *Engine) SetStartIndex(index int) {
	e.params.StartIndex = index
}

// SetLineCount updates the number of lines to draw.
func (e *Engine) SetLineCount(count int) {
	if count < 0 {
		count = 0
	}
	if count > e.params.PointCount {
		count = e.params.PointCount
	}
	e.params.LineCount = count
}

// SetLineAll toggles the draw-all mode for lines.
func (e *Engine) SetLineAll(all bool) {
	if all {
		e.params.LineCount = -1
		return
	}
	if e.params.LineCount < 0 {
		e.params.LineCount = e.params.PointCount
	}
}

// SetShowCircle toggles the circle outline.
func (e *Engine) SetShowCircle(show bool) {
	e.params.ShowCircle = show
}

// SetShowPoints toggles point rendering.
func (e *Engine) SetShowPoints(show bool) {
	e.params.ShowPoints = show
}

// SetShowLabels toggles label rendering.
func (e *Engine) SetShowLabels(show bool) {
	e.params.ShowLabels = show
}

// SetLabelStep updates the label step size.
func (e *Engine) SetLabelStep(step int) {
	if step < 1 {
		step = 1
	}
	e.params.LabelStep = step
}

// SetLineWidth updates the line width in CSS pixels.
func (e *Engine) SetLineWidth(width float64) {
	if width <= 0 {
		width = 1
	}
	e.params.LineWidth = width
}

// SetPointRadius updates the point radius in CSS pixels.
func (e *Engine) SetPointRadius(radius float64) {
	if radius < 0 {
		radius = 0
	}
	e.params.PointRadius = radius
}

// SetBackgroundColor updates the background color.
func (e *Engine) SetBackgroundColor(color string) {
	e.params.Colors.Background = color
}

// SetLineColor updates the line color.
func (e *Engine) SetLineColor(color string) {
	e.params.Colors.Line = color
}

// SetCircleColor updates the circle color.
func (e *Engine) SetCircleColor(color string) {
	e.params.Colors.Circle = color
}

// SetPointColor updates the point color.
func (e *Engine) SetPointColor(color string) {
	e.params.Colors.Point = color
}

// SetLabelColor updates the label color.
func (e *Engine) SetLabelColor(color string) {
	e.params.Colors.Label = color
}

// SetLineAnimation updates the line animation settings.
func (e *Engine) SetLineAnimation(settings AnimationSettings) {
	e.applyAnimationSettings(&e.animations.Lines, settings)
}

// SetMultiplierAnimation updates the multiplier animation settings.
func (e *Engine) SetMultiplierAnimation(settings AnimationSettings) {
	e.applyAnimationSettings(&e.animations.Multiplier, settings)
}

// SetPointAnimation updates the points animation settings.
func (e *Engine) SetPointAnimation(settings AnimationSettings) {
	e.applyAnimationSettings(&e.animations.Points, settings)
}

func (e *Engine) applyAnimationSettings(animation *Animation, settings AnimationSettings) {
	wasEnabled := animation.Settings.Enabled
	animation.Settings = settings
	if settings.Speed < 0 {
		animation.Settings.Speed = math.Abs(settings.Speed)
	}
	if !wasEnabled && settings.Enabled {
		animation.Value = settings.Start
		animation.Forward = true
		return
	}
	minV, maxV := ordered(settings.Start, settings.End)
	if animation.Value < minV || animation.Value > maxV {
		animation.Value = settings.Start
		animation.Forward = true
	}
}

func ordered(a, b float64) (float64, float64) {
	if a <= b {
		return a, b
	}
	return b, a
}

// Advance steps the animation forward and returns the new value.
func (a *Animation) Advance(dt float64) float64 {
	settings := a.Settings
	if !settings.Enabled || settings.Speed == 0 {
		return a.Value
	}
	if settings.Start == settings.End {
		a.Value = settings.Start
		return a.Value
	}

	minV, maxV := ordered(settings.Start, settings.End)
	direction := 1.0
	if settings.Start <= settings.End {
		if !a.Forward {
			direction = -1
		}
	} else {
		if a.Forward {
			direction = -1
		}
	}

	a.Value += direction * settings.Speed * dt

	if a.Value > maxV {
		a.handleBoundary(maxV, minV)
	} else if a.Value < minV {
		a.handleBoundary(minV, maxV)
	}

	return a.Value
}

func (a *Animation) handleBoundary(boundary, opposite float64) {
	settings := a.Settings
	if settings.PingPong {
		a.Value = boundary
		a.Forward = !a.Forward
		return
	}
	if settings.Loop {
		a.Value = opposite
		return
	}
	a.Value = boundary
	a.Settings.Enabled = false
}
