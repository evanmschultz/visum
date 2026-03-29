package app

import (
	"strings"
	"testing"

	"github.com/evanschultz/visum/internal/core"
)

func TestSVGExporterBasic(t *testing.T) {
	params := core.DefaultParams()
	params.PointCount = 6
	params.LineCount = 6
	params.LineWidth = 1.2
	params.PointRadius = 1.5
	params.ShowCircle = true
	params.ShowPoints = true
	params.ShowLabels = true
	params.LabelStep = 2

	exporter := NewSVGExporter()
	svg := exporter.Export(params, core.Size{Width: 200, Height: 200})

	if !strings.HasPrefix(svg, "<svg") {
		t.Fatalf("expected svg root, got %q", svg[:10])
	}
	if strings.Count(svg, "<line ") != 6 {
		t.Fatalf("expected 6 line elements")
	}
	circleCount := strings.Count(svg, "<circle ")
	if circleCount != 7 {
		t.Fatalf("expected 7 circles (points + outer), got %d", circleCount)
	}
	if strings.Count(svg, "<text ") == 0 {
		t.Fatalf("expected text labels to be rendered")
	}
}

func TestSVGExporterDefaultsSize(t *testing.T) {
	params := core.DefaultParams()
	exporter := NewSVGExporter()
	svg := exporter.Export(params, core.Size{})

	if !strings.Contains(svg, "width=\"800.00\"") || !strings.Contains(svg, "height=\"600.00\"") {
		t.Fatalf("expected default size in svg")
	}
}

func TestSVGExporterReadout(t *testing.T) {
	params := core.DefaultParams()
	params.Multiplier = 7.25
	exporter := NewSVGExporter()
	svg := exporter.ExportWithReadout(params, core.Size{Width: 300, Height: 200}, true)

	if !strings.Contains(svg, "k=7.25") {
		t.Fatalf("expected k readout in svg")
	}
}
