// pcb makes SVG files for conversion to Gerber to fabricate a PC Board.
package main

import (
	"flag"
	"log"
	"math"

	"github.com/gmlewis/sdfx/render"
	. "github.com/gmlewis/sdfx/sdf"
)

const (
	start1 = 2 * math.Pi
)

var (
	n     = flag.Int("n", 50, "Number of coils")
	round = flag.Float64("round", 0.25*math.Pi, "Round radius for spiral")
	svg   = flag.String("svg", "pcb", "Output SVG root filename")
)

func main() {
	flag.Parse()

	end1 := 2.0 * math.Pi * float64(*n)
	s := Spiral2D(start1, end1, *round)

	if *svg != "" {
		render.ToSVG(s, *svg+".svg", render.NewMarchingSquaresQuadtree(600))
	}

	log.Printf("Done.")
}
