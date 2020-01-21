// Computes an SVG rendering of a 3-D surface function
package main

import (
	"fmt"
	"io"
	"math"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange //pixels per x or y uintx
	zscale        = height * 0.4        // pixels pre z uint
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func svg(out io.Writer, color string) {
	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, _, ok := corner(i+1, j, color)
			if !ok {
				continue
			}
			bx, by, _, ok := corner(i, j, color)
			if !ok {
				continue
			}
			cx, cy, _, ok := corner(i, j+1, color)
			if !ok {
				continue
			}
			dx, dy, color, ok := corner(i+1, j+1, color)
			if !ok {
				continue
			}
			fmt.Fprintf(out, "<polygon style='stroke: %s' points='%g,%g %g,%g %g,%g %g,%g'/>\n", color, ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintln(out, "</svg>")
}

func corner(i, j int, c string) (float64, float64, string, bool) {
	// find point (x,y) at corner of cell (i,j)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// compute surface height z
	z, ok := f(x, y)
	if !ok {
		return 0, 0, "", false
	}

	// project (x,y,z) isometrically onto 2-D SVG canvas (sx, sy)
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	color := c
	if z < 0 {
		color = "#0000ff"
	}
	return sx, sy, color, true
}

func f(x, y float64) (float64, bool) {
	r := math.Hypot(x, y) // distance from (0,0)

	if math.IsInf(math.Sin(r)/r, 0) {
		return 0, false
	}
	return math.Sin(r) / r, true
}

// go run *.go
