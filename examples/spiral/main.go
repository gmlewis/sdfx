// spiral generates a PNG/SVG/DXF of a spiral.
package main

import (
	"flag"
	"log"
	"math"

	. "github.com/gmlewis/sdfx/sdf"
)

var (
	start = flag.Float64("start", 0.0, "Start radius (and angle) in radians of spiral")
	end   = flag.Float64("end", 2*math.Pi, "End radius (and angle) in radians of spiral")
	round = flag.Float64("round", 0.0, "Round radius for spiral")
	size  = flag.Int("size", 800, "Size of output file (width and height)")
	out   = flag.String("out", "spiral.png", "Output PNG filename of spiral")
	svg   = flag.String("svg", "spiral.svg", "Output SVG filename of spiral")
	dxf   = flag.String("dxf", "", "Output DXF filename of spiral")
)

func main() {
	flag.Parse()

	s := Spiral2D(*start, *end, *round)

	if *out != "" {
		png, err := NewPNG(*out, s.BoundingBox(), V2i{*size, *size})
		if err != nil {
			log.Fatalf("NewPNG: %v", err)
		}
		png.RenderSDF2(s)
		if err := png.Save(); err != nil {
			log.Fatalf("Save: %v", err)
		}
	}

	if *dxf != "" {
		RenderDXF(s, *size, *dxf)
	}

	if *svg != "" {
		if err := RenderSVG(s, *size, *svg); err != nil {
			log.Fatalf("RenderSVG: %v", err)
		}
	}
}
