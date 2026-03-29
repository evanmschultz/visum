package app

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/evanschultz/visum/internal/core"
)

// SVGExporter renders the current frame geometry as an SVG document.
type SVGExporter struct{}

// NewSVGExporter returns a new SVG exporter.
func NewSVGExporter() *SVGExporter {
	return &SVGExporter{}
}

// Export converts the provided params and size into a standalone SVG string.
func (e *SVGExporter) Export(params core.Params, size core.Size) string {
	return e.ExportWithReadout(params, size, false)
}

// ExportWithReadout converts the provided params and size into a standalone SVG string,
// optionally including the multiplier readout in the lower-left corner.
func (e *SVGExporter) ExportWithReadout(params core.Params, size core.Size, includeReadout bool) string {
	p := core.NormalizeParams(params)
	if size.Width <= 0 || size.Height <= 0 {
		size = core.Size{Width: 800, Height: 600}
	}
	frame := core.BuildFrame(p, size)

	var b strings.Builder
	b.Grow(4096)
	fmt.Fprintf(&b, "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%s\" height=\"%s\" viewBox=\"0 0 %s %s\">", svgFloat(size.Width), svgFloat(size.Height), svgFloat(size.Width), svgFloat(size.Height))
	fmt.Fprintf(&b, "<rect width=\"100%%\" height=\"100%%\" fill=\"%s\"/>", p.Colors.Background)

	if len(frame.Lines) > 0 {
		fmt.Fprintf(&b, "<g fill=\"none\" stroke=\"%s\" stroke-width=\"%s\" stroke-linecap=\"round\">", p.Colors.Line, svgFloat(p.LineWidth))
		for _, line := range frame.Lines {
			fmt.Fprintf(&b, "<line x1=\"%s\" y1=\"%s\" x2=\"%s\" y2=\"%s\"/>", svgFloat(line.From.X), svgFloat(line.From.Y), svgFloat(line.To.X), svgFloat(line.To.Y))
		}
		b.WriteString("</g>")
	}

	if p.ShowCircle {
		fmt.Fprintf(&b, "<circle cx=\"%s\" cy=\"%s\" r=\"%s\" fill=\"none\" stroke=\"%s\" stroke-width=\"%s\"/>", svgFloat(frame.Circle.Center.X), svgFloat(frame.Circle.Center.Y), svgFloat(frame.Circle.Radius), p.Colors.Circle, svgFloat(p.LineWidth))
	}

	if p.ShowPoints {
		fmt.Fprintf(&b, "<g fill=\"%s\">", p.Colors.Point)
		for _, point := range frame.Points {
			fmt.Fprintf(&b, "<circle cx=\"%s\" cy=\"%s\" r=\"%s\"/>", svgFloat(point.X), svgFloat(point.Y), svgFloat(p.PointRadius))
		}
		b.WriteString("</g>")
	}

	if p.ShowLabels && len(frame.Labels) > 0 {
		fontSize := math.Max(10, frame.Circle.Radius*0.06)
		fmt.Fprintf(&b, "<g fill=\"%s\" font-family=\"Source Serif 4, Iowan Old Style, Palatino Linotype, serif\" font-size=\"%s\" font-weight=\"300\" text-anchor=\"middle\" dominant-baseline=\"middle\">", p.Colors.Label, svgFloat(fontSize))
		for _, label := range frame.Labels {
			fmt.Fprintf(&b, "<text x=\"%s\" y=\"%s\">%s</text>", svgFloat(label.Position.X), svgFloat(label.Position.Y), label.Text)
		}
		b.WriteString("</g>")
	}

	if includeReadout {
		readout := fmt.Sprintf("k=%s", formatReadout(p.Multiplier, size.Width))
		fontSize := math.Max(12, size.Width*0.02)
		x := 14.0
		y := size.Height - 14.0
		fmt.Fprintf(&b, "<text x=\"%s\" y=\"%s\" fill=\"%s\" font-family=\"Source Serif 4, Iowan Old Style, Palatino Linotype, serif\" font-size=\"%s\" font-weight=\"300\" text-anchor=\"start\" dominant-baseline=\"alphabetic\">%s</text>", svgFloat(x), svgFloat(y), p.Colors.Label, svgFloat(fontSize), readout)
	}

	b.WriteString("</svg>")
	return b.String()
}

func svgFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func formatReadout(value float64, width float64) string {
	precision := 3
	if width < 520 {
		precision = 2
	}
	if math.Abs(value) >= 100 {
		precision = 1
	}
	return strconv.FormatFloat(value, 'f', precision, 64)
}
