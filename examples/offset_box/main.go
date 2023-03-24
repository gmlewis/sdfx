//-----------------------------------------------------------------------------
/*

Build a box using offsets from a rectangular box.

TODO Add a retaining lip to the base or top so the lid stays in place.

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"

	"github.com/gmlewis/sdfx/render"
	"github.com/gmlewis/sdfx/sdf"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

const sizeX = 30.0
const sizeY = 40.0
const sizeZ = 30.0

const wallThickness = 3.0
const outerRadius = 6.0
const lidPosition = 0.75 // 0..1 position of lid on box

//-----------------------------------------------------------------------------

func box() error {

	if outerRadius < wallThickness {
		return sdf.ErrMsg("outerRadius < wallThickness")
	}

	innerOfs := outerRadius - wallThickness
	outerOfs := innerOfs + wallThickness

	if sizeX < outerOfs {
		return sdf.ErrMsg("sizeX < outerOfs")
	}
	if sizeY < outerOfs {
		return sdf.ErrMsg("sizeY < outerOfs")
	}
	if sizeZ < outerOfs {
		return sdf.ErrMsg("sizeZ < outerOfs")
	}

	baseBox, err := sdf.Box3D(v3.Vec{sizeX - outerOfs, sizeY - outerOfs, sizeZ - outerOfs}, 0)
	if err != nil {
		return err
	}

	innerBox := sdf.Offset3D(baseBox, innerOfs)
	outerBox := sdf.Offset3D(baseBox, outerOfs)
	box := sdf.Difference3D(outerBox, innerBox)

	lidZ := (lidPosition - 0.5) * sizeZ
	base := sdf.Cut3D(box, v3.Vec{0, 0, lidZ}, v3.Vec{0, 0, -1})
	top := sdf.Cut3D(box, v3.Vec{0, 0, lidZ}, v3.Vec{0, 0, 1})

	render.ToSTL(base, "base.stl", render.NewMarchingCubesOctree(300))
	render.ToSTL(top, "top.stl", render.NewMarchingCubesOctree(300))

	return nil
}

//-----------------------------------------------------------------------------

func main() {
	err := box()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

//-----------------------------------------------------------------------------
