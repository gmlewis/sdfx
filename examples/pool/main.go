//-----------------------------------------------------------------------------
/*

Pool Model

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/gmlewis/sdfx/render"
	"github.com/gmlewis/sdfx/sdf"
	v2 "github.com/gmlewis/sdfx/vec/v2"
)

//-----------------------------------------------------------------------------

const cubicInchesPerGallon = 231.0

// pool dimensions are in inches
const poolWidth = 234.0
const poolLength = 477.0

var poolDepth = []v2.Vec{
	{0.0, 43.0},
	{101.0, 46.0},
	{202.0, 58.0},
	{298.0, 83.0},
	{394.0, 96.0},
	{477.0, 96.0},
}

const vol = (7738.3005 * 1000.0) / cubicInchesPerGallon // gallons

//-----------------------------------------------------------------------------

func pool() (sdf.SDF3, error) {
	log.Printf("pool volume %f gallons\n", vol)
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.AddV2Set(poolDepth)
	p.Add(poolLength, 0)
	profile, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(profile, poolWidth), nil
}

//-----------------------------------------------------------------------------

func main() {
	pool, err := pool()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(pool, "pool1.stl", render.NewMarchingCubesOctree(300))
	//render.ToSTL(pool, 15, "pool2.stl", dc.NewDualContouringDefault())
}

//-----------------------------------------------------------------------------
