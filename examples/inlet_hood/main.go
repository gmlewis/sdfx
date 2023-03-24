//-----------------------------------------------------------------------------
/*

Inlet Masking Hood

As seen on various GDI engines:
Covers an Audi/VW inlet port while walnut blasting carbon deposits.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/gmlewis/sdfx/obj"
	"github.com/gmlewis/sdfx/render"
	"github.com/gmlewis/sdfx/sdf"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//const shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

var baseSize = v3.Vec{40, 60, 10} // 20
var portSize = v3.Vec{30, 50, 10} // 15

//-----------------------------------------------------------------------------

func outerBase() (sdf.SDF3, error) {

	trp := &obj.TruncRectPyramidParms{
		Size:        baseSize,
		BaseAngle:   sdf.DtoR(90.0 - 2.0),
		BaseRadius:  baseSize.X * 0.5,
		RoundRadius: 0,
	}

	return obj.TruncRectPyramid3D(trp)
}

func innerBase() (sdf.SDF3, error) {

	trp := &obj.TruncRectPyramidParms{
		Size:        portSize,
		BaseAngle:   sdf.DtoR(90.0 - 5.0),
		BaseRadius:  portSize.X * 0.5,
		RoundRadius: 0,
	}

	return obj.TruncRectPyramid3D(trp)
}

func hood() (sdf.SDF3, error) {

	ob, err := outerBase()
	if err != nil {
		return nil, err
	}

	ib, err := innerBase()
	if err != nil {
		return nil, err
	}

	return sdf.Difference3D(ob, ib), nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := hood()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(s, shrink), "hood.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
