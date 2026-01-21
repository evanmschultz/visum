package core

import (
	"math"
	"testing"
)

func TestPointsOnCircleCount(t *testing.T) {
	points := PointsOnCircle(12, 10, 0, Vec2{})
	if len(points) != 12 {
		t.Fatalf("expected 12 points, got %d", len(points))
	}
}

func TestTimesTableLinesCount(t *testing.T) {
	lines := TimesTableLines(10, 1, 0, Vec2{}, 2, 0, 5)
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}
}

func TestTimesTableLinesFirstLine(t *testing.T) {
	lines := TimesTableLines(10, 1, 0, Vec2{}, 2, 0, 1)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	line := lines[0]
	if !almostEqual(line.From.X, 0) || !almostEqual(line.From.Y, -1) {
		t.Fatalf("unexpected line start: %+v", line.From)
	}
	if !almostEqual(line.To.X, 0) || !almostEqual(line.To.Y, -1) {
		t.Fatalf("unexpected line end: %+v", line.To)
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}
