//-----------------------------------------------------------------------------
/*

Pillar Holder

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/gmlewis/sdfx/render"
	"github.com/gmlewis/sdfx/sdf"
	v2 "github.com/gmlewis/sdfx/vec/v2"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

var wallThickness = 2.5
var wallHeight = 15.0
var pillarWidth = 33.0
var pillarRadius = 4.0
var feetWidth = 6.0
var baseThickness = 3.0

//-----------------------------------------------------------------------------

func base() sdf.SDF3 {
	w := pillarWidth + 2.0*(feetWidth+wallThickness)
	h := pillarWidth + 2.0*wallThickness
	r := pillarRadius + wallThickness
	base2d := sdf.Box2D(v2.Vec{w, h}, r)
	return sdf.Extrude3D(base2d, baseThickness)
}

func wall(w, r float64) sdf.SDF3 {
	base := sdf.Box2D(v2.Vec{w, w}, r)
	s := sdf.Extrude3D(base, wallHeight)
	ofs := 0.5 * (wallHeight - baseThickness)
	return sdf.Transform3D(s, sdf.Translate3d(v3.Vec{0, 0, ofs}))
}

func holder() sdf.SDF3 {
	base := base()
	outer := wall(pillarWidth+2.0*wallThickness, pillarRadius+wallThickness)
	inner := wall(pillarWidth, pillarRadius)
	return sdf.Difference3D(sdf.Union3D(base, outer), inner)
}

//-----------------------------------------------------------------------------

func main() {
	s := holder()
	s = sdf.ScaleUniform3D(s, shrink)
	render.ToSTL(s, "holder.stl", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
