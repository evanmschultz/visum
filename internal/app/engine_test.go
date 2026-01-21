package app

import (
	"math"
	"testing"

	"github.com/evanschultz/visum/internal/core"
)

func TestEngineStepMultiplier(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetStepTarget(StepMultiplier)
	engine.SetStepAmount(0.5)

	before := engine.Snapshot().Params.Multiplier
	engine.Step(1)
	after := engine.Snapshot().Params.Multiplier

	if !almostEqual(after, before+0.5) {
		t.Fatalf("expected multiplier %.2f, got %.2f", before+0.5, after)
	}
}

func TestEngineStepLines(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetLineCount(20)
	engine.SetStepTarget(StepLines)
	engine.SetStepAmount(2)

	engine.Step(-1)
	if engine.Snapshot().Params.LineCount != 18 {
		t.Fatalf("expected line count 18, got %d", engine.Snapshot().Params.LineCount)
	}
}

func TestEngineUpdateLineAnimation(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetLineAnimation(AnimationSettings{
		Enabled: true,
		Start:   0,
		End:     10,
		Speed:   10,
	})
	engine.SetRunning(true)

	engine.Update(1)
	if engine.Snapshot().Params.LineCount != 10 {
		t.Fatalf("expected line count 10, got %d", engine.Snapshot().Params.LineCount)
	}
}

func TestAnimationPingPong(t *testing.T) {
	anim := Animation{Settings: AnimationSettings{Enabled: true, Start: 0, End: 1, Speed: 2, PingPong: true}, Value: 0, Forward: true}
	anim.Advance(1)
	if !almostEqual(anim.Value, 1) || anim.Forward {
		t.Fatalf("expected pingpong to hit end and reverse, got %.2f forward=%v", anim.Value, anim.Forward)
	}

	anim.Advance(0.5)
	if anim.Value >= 1 {
		t.Fatalf("expected value to move backward, got %.2f", anim.Value)
	}
}

func TestSetLineAll(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetLineAll(true)
	if engine.Snapshot().Params.LineCount >= 0 {
		t.Fatalf("expected line count to be -1 when all lines enabled")
	}

	engine.SetLineAll(false)
	if engine.Snapshot().Params.LineCount != engine.Snapshot().Params.PointCount {
		t.Fatalf("expected line count to reset to point count")
	}
}

func TestSetPointCountClamp(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetPointCount(1)
	if engine.Snapshot().Params.PointCount != 2 {
		t.Fatalf("expected point count to clamp to 2")
	}

	engine.SetPointCount(5000)
	if engine.Snapshot().Params.PointCount != 4000 {
		t.Fatalf("expected point count to clamp to 4000")
	}
}

func TestResetDefaults(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetMultiplier(9)
	engine.SetPointCount(1000)
	engine.Reset(core.DefaultParams())

	snapshot := engine.Snapshot()
	if snapshot.Params.PointCount != core.DefaultParams().PointCount {
		t.Fatalf("expected point count reset to default")
	}
	if !almostEqual(snapshot.Params.Multiplier, core.DefaultParams().Multiplier) {
		t.Fatalf("expected multiplier reset to default")
	}
}

func TestEngineSetters(t *testing.T) {
	engine := NewEngine(core.DefaultParams())

	engine.SetRotationDeg(30)
	engine.SetStartIndex(5)
	engine.SetShowCircle(false)
	engine.SetShowPoints(true)
	engine.SetShowLabels(true)
	engine.SetLabelStep(0)
	engine.SetLineWidth(0)
	engine.SetPointRadius(-1)

	engine.SetBackgroundColor("#000000")
	engine.SetLineColor("#111111")
	engine.SetCircleColor("#222222")
	engine.SetPointColor("#333333")
	engine.SetLabelColor("#444444")

	params := engine.Snapshot().Params
	if !almostEqual(params.RotationDeg, 30) {
		t.Fatalf("expected rotation 30, got %.2f", params.RotationDeg)
	}
	if params.StartIndex != 5 {
		t.Fatalf("expected start index 5, got %d", params.StartIndex)
	}
	if params.ShowCircle {
		t.Fatalf("expected show circle false")
	}
	if !params.ShowPoints || !params.ShowLabels {
		t.Fatalf("expected points and labels to be enabled")
	}
	if params.LabelStep != 1 {
		t.Fatalf("expected label step to clamp to 1")
	}
	if params.LineWidth != 1 {
		t.Fatalf("expected line width to clamp to 1")
	}
	if params.PointRadius != 0 {
		t.Fatalf("expected point radius to clamp to 0")
	}
	if params.Colors.Background != "#000000" || params.Colors.Line != "#111111" || params.Colors.Circle != "#222222" || params.Colors.Point != "#333333" || params.Colors.Label != "#444444" {
		t.Fatalf("unexpected colors: %+v", params.Colors)
	}
}

func TestEngineStepPoints(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetStepTarget(StepPoints)
	engine.SetStepAmount(3.4)

	engine.Step(1)
	if engine.Snapshot().Params.PointCount != core.DefaultParams().PointCount+3 {
		t.Fatalf("expected points to increase by 3")
	}
	engine.Step(-1)
	if engine.Snapshot().Params.PointCount != core.DefaultParams().PointCount {
		t.Fatalf("expected points to return to default")
	}
}

func TestUpdateNoRun(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetLineAnimation(AnimationSettings{Enabled: true, Start: 0, End: 5, Speed: 10})
	before := engine.Snapshot().Params.LineCount
	engine.Update(1)

	if engine.Snapshot().Params.LineCount != before {
		t.Fatalf("expected line count to remain unchanged when not running")
	}
}

func TestAnimationLoop(t *testing.T) {
	anim := Animation{Settings: AnimationSettings{Enabled: true, Start: 0, End: 1, Speed: 3, Loop: true}, Value: 0, Forward: true}
	anim.Advance(1)
	if anim.Value != 0 {
		t.Fatalf("expected loop to wrap to start, got %.2f", anim.Value)
	}
}

func TestAnimationReverseRange(t *testing.T) {
	anim := Animation{Settings: AnimationSettings{Enabled: true, Start: 5, End: 1, Speed: 2, PingPong: true}, Value: 5, Forward: true}
	anim.Advance(1)
	if anim.Value > 5 {
		t.Fatalf("expected value to move toward lower bound, got %.2f", anim.Value)
	}
}

func TestAnimationNegativeSpeed(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetMultiplierAnimation(AnimationSettings{Enabled: true, Start: 1, End: 2, Speed: -0.5})
	settings := engine.Snapshot().Animations.Multiplier.Settings
	if settings.Speed <= 0 {
		t.Fatalf("expected negative speed to be normalized")
	}
}

func TestEngineReverse(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetMultiplierAnimation(AnimationSettings{Enabled: true, Start: 0, End: 10, Speed: 1})
	engine.SetRunning(true)
	engine.Update(1)
	forwardValue := engine.Snapshot().Params.Multiplier

	engine.SetReverse(true)
	engine.Update(1)
	reversedValue := engine.Snapshot().Params.Multiplier

	if reversedValue >= forwardValue {
		t.Fatalf("expected reversed animation to move backward")
	}
}

func TestStepLinesFromAll(t *testing.T) {
	engine := NewEngine(core.DefaultParams())
	engine.SetLineAll(true)
	engine.SetStepTarget(StepLines)
	engine.SetStepAmount(5)

	engine.Step(-1)
	if engine.Snapshot().Params.LineCount != engine.Snapshot().Params.PointCount-5 {
		t.Fatalf("expected line count to step down from all")
	}
}


func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}
