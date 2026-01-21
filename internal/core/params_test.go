package core

import "testing"

func TestNormalizeParams(t *testing.T) {
	params := NormalizeParams(Params{
		PointCount: 0,
		LabelStep:  0,
		LineWidth:  0,
		PointRadius: -1,
		StartIndex: -1,
	})

	if params.PointCount != 2 {
		t.Fatalf("expected point count to clamp to 2")
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
	if params.StartIndex != 1 {
		t.Fatalf("expected start index to wrap into range")
	}
}

func TestBuildFrameLabels(t *testing.T) {
	params := DefaultParams()
	params.PointCount = 10
	params.ShowLabels = true
	params.LabelStep = 2

	frame := BuildFrame(params, Size{Width: 200, Height: 200})
	if len(frame.Labels) != 5 {
		t.Fatalf("expected 5 labels, got %d", len(frame.Labels))
	}
}
