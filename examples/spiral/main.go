//-----------------------------------------------------------------------------
/*

Spirals

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/gmlewis/sdfx/render"
	"github.com/gmlewis/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {
	s, err := sdf.ArcSpiral2D(1.0, 20.0, 0.25*sdf.Pi, 8*sdf.Tau, 1.0)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	render.ToDXF(s, "spiral.dxf", render.NewMarchingSquaresQuadtree(400))
}

//-----------------------------------------------------------------------------
