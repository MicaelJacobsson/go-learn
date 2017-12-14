// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		surface(w)
	}
	http.HandleFunc("/", handler)
	//!-http
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}

func surface(out io.Writer) {

	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			all_ok := true
			var colour uint32
			ok, ax, ay, c := corner(i+1, j)
			all_ok = all_ok && ok
			colour += c
			ok, bx, by, c := corner(i, j)
			all_ok = all_ok && ok

			ok, cx, cy, c := corner(i, j+1)
			all_ok = all_ok && ok

			ok, dx, dy, c := corner(i+1, j+1)
			all_ok = all_ok && ok

			if !all_ok {
				continue
			}
			fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='#%06x'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, colour)
		}
	}
	fmt.Fprintf(out, "</svg>")
}

func corner(i, j int) (bool, float64, float64, uint32) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	// Colour ff0000 to 0000ff
	var red, blue, colour uint32
	red = uint32(256.0 * sy / float64(height))
	blue = uint32(256.0 * (1 - (sy / float64(height))))
	colour = ((red << 16) & 0xff0000) + (blue & 0xff)

	return !math.IsInf(z, 0), sx, sy, colour
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

//!-
