package core

// Size represents a 2D extent in CSS pixels.
type Size struct {
	Width  float64
	Height float64
}

// Vec2 is a 2D vector or point.
type Vec2 struct {
	X float64
	Y float64
}

// Line represents a line segment in 2D space.
type Line struct {
	From Vec2
	To   Vec2
}

// Label represents a text label placed on the canvas.
type Label struct {
	Position Vec2
	Text     string
}

// Circle represents a circle to draw on the canvas.
type Circle struct {
	Center Vec2
	Radius float64
}

// Colors defines CSS color strings used by the renderer.
type Colors struct {
	Background string
	Line       string
	Circle     string
	Point      string
	Label      string
}

// Params defines the user-controlled parameters for rendering.
type Params struct {
	PointCount int
	Multiplier float64
	RotationDeg float64
	StartIndex int
	// LineCount is the number of lines to draw. Use -1 to draw all lines.
	LineCount int

	ShowCircle bool
	ShowPoints bool
	ShowLabels bool
	LabelStep  int

	LineWidth  float64
	PointRadius float64

	Colors Colors
}

// Frame is the fully resolved geometry for rendering.
type Frame struct {
	Circle Circle
	Lines  []Line
	Points []Vec2
	Labels []Label
}

// DefaultParams returns a baseline configuration for the app.
func DefaultParams() Params {
	return Params{
		PointCount: 200,
		Multiplier: 2,
		RotationDeg: 0,
		StartIndex: 0,
		LineCount: -1,
		ShowCircle: true,
		ShowPoints: false,
		ShowLabels: false,
		LabelStep: 10,
		LineWidth: 1.0,
		PointRadius: 1.5,
		Colors: Colors{
			Background: "#0b0f1a",
			Line: "#3a6ff7",
			Circle: "#ff4d4d",
			Point: "#f2f4f8",
			Label: "#f2f4f8",
		},
	}
}
