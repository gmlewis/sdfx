//-----------------------------------------------------------------------------
/*

Pipe Flange with a Square base

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/gmlewis/sdfx/obj"
	"github.com/gmlewis/sdfx/render"
	"github.com/gmlewis/sdfx/sdf"
	"github.com/gmlewis/sdfx/vec/conv"
	v2 "github.com/gmlewis/sdfx/vec/v2"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

const pipeClearance = 1.01                 // ID of pipe holder slightly larger to allow a slip fit
const pipeDiameter = 48.45 * pipeClearance // OD of pipe to be fitted
var baseSize = v2.Vec{77.0, 77.0}          // size of rectangular base
const baseThickness = 3.0                  // base thickness
const pipeWall = 3.0                       // pipe holder wall thickness
const pipeLength = 30.0                    // length of pipe holder (from bottom)
var pipeOffset = v2.Vec{0, 0}              // offset of pipe holder from base center

const pipeRadius = 0.5 * pipeDiameter
const pipeFillet = 0.95 * pipeWall

//-----------------------------------------------------------------------------

func flange() (sdf.SDF3, error) {

	// base
	pp := &obj.PanelParms{
		Size:         baseSize,
		CornerRadius: 18.0,
		HoleDiameter: 3.5, // #6 screw
		HoleMargin:   [4]float64{12.0, 12.0, 12.0, 12.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}
	panel, err := obj.Panel2D(pp)
	if err != nil {
		return nil, err
	}
	base := sdf.Extrude3D(panel, 2.0*baseThickness)

	// outer pipe
	outerPipe, _ := sdf.Cylinder3D(2.0*pipeLength, pipeRadius+pipeWall, 0.0)
	outerPipe = sdf.Transform3D(outerPipe, sdf.Translate3d(conv.V2ToV3(pipeOffset, 0)))

	// inner pipe
	innerPipe, _ := sdf.Cylinder3D(2.0*pipeLength, pipeRadius, 0.0)
	innerPipe = sdf.Transform3D(innerPipe, sdf.Translate3d(conv.V2ToV3(pipeOffset, 0)))

	// combine the outer pipe and base (with a fillet)
	s0 := sdf.Union3D(base, outerPipe)
	s0.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(pipeFillet))

	// cut the through hole
	s := sdf.Difference3D(s0, innerPipe)

	// return the upper half
	return sdf.Cut3D(s, v3.Vec{0, 0, 0}, v3.Vec{0, 0, 1}), nil
}

//-----------------------------------------------------------------------------

func main() {
	flange, err := flange()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(sdf.ScaleUniform3D(flange, shrink), "flange.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
