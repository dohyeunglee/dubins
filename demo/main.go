package main

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"image/color"

	"github.com/dohyeunglee/dubins"
)

const (
	stepSize      = 0.5
	turningRadius = 1.0
)

var (
	examples = []struct {
		Name    string
		Start   dubins.State
		Goal    dubins.State
		AxisMin [2]float64
		AxisMax [2]float64
	}{
		{
			Name:  "(0, 0, 0) -> (-4, 0, 0)",
			Start: dubins.State{X: 0, Y: 0, Yaw: 0}, Goal: dubins.State{X: -4, Y: 0, Yaw: 0},
			AxisMin: [2]float64{-6, -0.5},
			AxisMax: [2]float64{2, 2.5},
		},
		{
			Name:  "(0, 0, 0) -> (4, 4, 0)",
			Start: dubins.State{X: 0, Y: 0, Yaw: 0}, Goal: dubins.State{X: 4, Y: 4, Yaw: 0},
			AxisMin: [2]float64{0, 0},
			AxisMax: [2]float64{5, 4},
		},
		{
			Name:  "(0, 0, 0) -> (4, -4, 0)",
			Start: dubins.State{X: 0, Y: 0, Yaw: 0}, Goal: dubins.State{X: 4, Y: -4, Yaw: 0},
			AxisMin: [2]float64{0, -4},
			AxisMax: [2]float64{5, 0},
		},
		{
			Name:  "(0, 0, 0) -> (-4, 4, 0)",
			Start: dubins.State{X: 0, Y: 0, Yaw: 0}, Goal: dubins.State{X: -4, Y: 4, Yaw: 0},
			AxisMin: [2]float64{-6, 0},
			AxisMax: [2]float64{2, 4},
		},
		{
			Name:  "(0, 0, 0) -> (-4, -4, 0)",
			Start: dubins.State{X: 0, Y: 0, Yaw: 0}, Goal: dubins.State{X: -4, Y: -4, Yaw: 0},
			AxisMin: [2]float64{-6, -4},
			AxisMax: [2]float64{2, 0},
		},
		{
			Name:  "(0, 0, 0) -> (4, 4, pi)",
			Start: dubins.State{X: 0, Y: 0, Yaw: 0}, Goal: dubins.State{X: 4, Y: 4, Yaw: math.Pi},
			AxisMin: [2]float64{0, 0},
			AxisMax: [2]float64{6, 4},
		},
		{
			Name:  "(0, 0, 0) -> (4, -4, pi)",
			Start: dubins.State{X: 0, Y: 0, Yaw: 0}, Goal: dubins.State{X: 4, Y: -4, Yaw: math.Pi},
			AxisMin: [2]float64{0, -4},
			AxisMax: [2]float64{6, 0},
		},
		{
			Name:  "(0, 0, 0) -> (0.5, 0, pi)",
			Start: dubins.State{X: 0, Y: 0, Yaw: 0}, Goal: dubins.State{X: 0.5, Y: 0, Yaw: math.Pi},
			AxisMin: [2]float64{-1, -3},
			AxisMax: [2]float64{5, 3},
		},
		{
			Name:  "(0, 0, pi/4) -> (4, 4, pi/4)",
			Start: dubins.State{X: 0, Y: 0, Yaw: math.Pi / 4}, Goal: dubins.State{X: 4, Y: 4, Yaw: math.Pi / 4},
			AxisMin: [2]float64{0, 0},
			AxisMax: [2]float64{5, 5},
		},
		{
			Name:  "(4, 4, 4/pi) -> (0, 0, pi/4)",
			Start: dubins.State{X: 4, Y: 4, Yaw: math.Pi / 4}, Goal: dubins.State{X: 0, Y: 0, Yaw: math.Pi / 4},
			AxisMin: [2]float64{-3, -1},
			AxisMax: [2]float64{5, 7},
		},
	}
)

func getPoint(centerX float64, centerY float64, radius float64, orin float64) (x float64, y float64) {
	x = centerX + radius*math.Cos(orin)
	y = centerY + radius*math.Sin(orin)
	return
}

func plotCar(p *plot.Plot, state dubins.State) {
	aX, aY := getPoint(state.X, state.Y, stepSize, state.Yaw)
	bX, bY := getPoint(state.X, state.Y, stepSize/2, state.Yaw+150.0/180.0*math.Pi)
	cX, cY := getPoint(state.X, state.Y, stepSize/2, state.Yaw-150.0/180.0*math.Pi)

	points := plotter.XYs{
		{X: aX, Y: aY},
		{X: bX, Y: bY},
		{X: cX, Y: cY},
		{X: aX, Y: aY},
	}

	line, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}

	line.Color = color.RGBA{G: 255, A: 255}
	line.Width = vg.Points(1.5)

	p.Add(line)
}

func plotPath(p *plot.Plot, states []dubins.State) {
	points := make(plotter.XYs, 0, len(states))
	for _, state := range states {
		points = append(points, plotter.XY{X: state.X, Y: state.Y})
	}

	lpLine, lpPoints, err := plotter.NewLinePoints(points)
	if err != nil {
		panic(err)
	}

	lpLine.Color = color.RGBA{B: 255, A: 255}
	lpPoints.Color = color.RGBA{R: 255, A: 255}
	lpPoints.Shape = draw.CircleGlyph{}

	p.Add(lpLine, lpPoints)
}

func main() {
	for i, example := range examples {
		path, ok := dubins.MinLengthPath(example.Start, example.Goal, turningRadius)
		if !ok {
			panic("MinLengthPath fail")
		}

		interpolated := path.Interpolate(stepSize)

		p := plot.New()
		p.Title.Text = fmt.Sprintf("DubinsPath Demo: %s\nLength: %.2f, Type: %s", example.Name, path.Length(), path.PathType())
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		p.Add(plotter.NewGrid())
		p.X.Min = example.AxisMin[0]
		p.Y.Min = example.AxisMin[1]
		p.X.Max = example.AxisMax[0]
		p.Y.Max = example.AxisMax[1]

		plotCar(p, example.Start)
		plotCar(p, example.Goal)

		plotPath(p, interpolated)

		if err := p.Save(4*vg.Inch, 4*vg.Inch, fmt.Sprintf("demo_path%d.png", i)); err != nil {
			panic(err)
		}
	}
}
