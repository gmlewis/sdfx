// pcb makes SVG files for conversion to Gerber to fabricate a PC Board.
package main

import (
	"flag"
	"log"
	"math"

	. "github.com/gmlewis/sdfx/sdf"
)

const (
	start1 = 2 * math.Pi
)

var (
	n         = flag.Int("n", 50, "Number of coils")
	round     = flag.Float64("round", 0.25*math.Pi, "Round radius for spiral")
	size      = flag.Int("size", 800, "Size of output file (width and height)")
	svg       = flag.String("svg", "pcb", "Output SVG root filename")
	lineStyle = flag.String("line_style", "fill:none;stroke:black;stroke-width:0.1", "SVG line style")
)

func main() {
	flag.Parse()

	end1 := 2.0 * math.Pi * float64(*n)
	s := Spiral2D(start1, end1, *round)

	if *svg != "" {
		if err := RenderSVG(s, *size, *svg+".svg", *lineStyle); err != nil {
			log.Fatalf("RenderSVG: %v", err)
		}

		if err := RenderSVGSlow(s, *size, *svg+"_slow.svg", *lineStyle); err != nil {
			log.Fatalf("RenderSVG_Slow: %v", err)
		}
	}
}
