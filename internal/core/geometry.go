package core

import (
	"fmt"
	"math"
)

// BuildFrame converts Params into concrete geometry for rendering.
func BuildFrame(params Params, size Size) Frame {
	p := NormalizeParams(params)
	if size.Width <= 0 || size.Height <= 0 {
		return Frame{}
	}

	center := Vec2{X: size.Width / 2, Y: size.Height / 2}
	radius := math.Min(size.Width, size.Height) * 0.42
	rotation := degToRad(p.RotationDeg)

	points := PointsOnCircle(p.PointCount, radius, rotation, center)
	lineCount := p.LineCount
	if lineCount < 0 || lineCount > p.PointCount {
		lineCount = p.PointCount
	}
	lines := TimesTableLines(p.PointCount, radius, rotation, center, p.Multiplier, p.StartIndex, lineCount)

	var labels []Label
	if p.ShowLabels {
		labelRadius := radius * 1.08
		labels = LabelsOnCircle(p.PointCount, labelRadius, rotation, center, p.LabelStep)
	}

	return Frame{
		Circle: Circle{Center: center, Radius: radius},
		Lines:  lines,
		Points: points,
		Labels: labels,
	}
}

// NormalizeParams clamps parameters to safe, usable ranges.
func NormalizeParams(params Params) Params {
	p := params
	if p.PointCount < 2 {
		p.PointCount = 2
	}
	if p.LabelStep < 1 {
		p.LabelStep = 1
	}
	if p.LineWidth <= 0 {
		p.LineWidth = 1
	}
	if p.PointRadius < 0 {
		p.PointRadius = 0
	}
	p.StartIndex = modInt(p.StartIndex, p.PointCount)
	return p
}

// PointsOnCircle returns evenly spaced points on a circle.
func PointsOnCircle(count int, radius float64, rotation float64, center Vec2) []Vec2 {
	if count < 1 {
		return nil
	}
	points := make([]Vec2, 0, count)
	baseAngle := -math.Pi/2 + rotation
	step := (2 * math.Pi) / float64(count)
	for i := 0; i < count; i++ {
		angle := baseAngle + step*float64(i)
		points = append(points, PointOnCircle(radius, angle, center))
	}
	return points
}

// LabelsOnCircle returns labels placed around a circle.
func LabelsOnCircle(count int, radius float64, rotation float64, center Vec2, step int) []Label {
	if count < 1 || step < 1 {
		return nil
	}
	labels := make([]Label, 0, count/step+1)
	baseAngle := -math.Pi/2 + rotation
	angleStep := (2 * math.Pi) / float64(count)
	for i := 0; i < count; i++ {
		if i%step != 0 {
			continue
		}
		angle := baseAngle + angleStep*float64(i)
		labels = append(labels, Label{
			Position: PointOnCircle(radius, angle, center),
			Text:     fmt.Sprintf("%d", i),
		})
	}
	return labels
}

// TimesTableLines returns line segments for the times-table mapping.
func TimesTableLines(count int, radius float64, rotation float64, center Vec2, multiplier float64, startIndex, lineCount int) []Line {
	if count < 2 || lineCount <= 0 {
		return nil
	}

	lines := make([]Line, 0, lineCount)
	baseAngle := -math.Pi/2 + rotation
	step := (2 * math.Pi) / float64(count)

	for i := 0; i < lineCount; i++ {
		index := modInt(startIndex+i, count)
		sourceAngle := baseAngle + step*float64(index)

		targetIndex := math.Mod(float64(index)*multiplier, float64(count))
		if targetIndex < 0 {
			targetIndex += float64(count)
		}
		targetAngle := baseAngle + step*targetIndex

		lines = append(lines, Line{
			From: PointOnCircle(radius, sourceAngle, center),
			To:   PointOnCircle(radius, targetAngle, center),
		})
	}

	return lines
}

// PointOnCircle returns a point on a circle at the given angle in radians.
func PointOnCircle(radius, angle float64, center Vec2) Vec2 {
	return Vec2{
		X: center.X + radius*math.Cos(angle),
		Y: center.Y + radius*math.Sin(angle),
	}
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

func modInt(value, mod int) int {
	if mod == 0 {
		return 0
	}
	result := value % mod
	if result < 0 {
		result += mod
	}
	return result
}
