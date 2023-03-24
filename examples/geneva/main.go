//-----------------------------------------------------------------------------
/*

3D printable geneva drive mechanism

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

var k0 = obj.GenevaParms{
	NumSectors:     6,
	CenterDistance: 50.0,
	DriverRadius:   20.0,
	DrivenRadius:   40.0,
	PinRadius:      2.5,
	Clearance:      0.1,
}

var k1 = obj.GenevaParms{
	NumSectors:     10,
	CenterDistance: 45.0,
	DriverRadius:   12.0,
	DrivenRadius:   45.0,
	PinRadius:      2.0,
	Clearance:      0.1,
}

func main() {

	k := k0

	sDriver, sDriven, err := obj.Geneva2D(&k)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	wheelHeight := 5.0                 // height of wheels
	holeRadius := 3.25                 // radius of center hole
	hubRadius := 10.0                  // hub radius for driven wheel
	baseRadius := 1.5 * k.DriverRadius // radius of base for driver wheel

	// extrude the driver wheel
	driver3d := sdf.Extrude3D(sDriver, wheelHeight)
	driver3d = sdf.Transform3D(driver3d, sdf.Translate3d(v3.Vec{0, 0, wheelHeight / 2}))
	// add a base
	base3d, _ := sdf.Cylinder3D(wheelHeight, baseRadius, 0)
	base3d = sdf.Transform3D(base3d, sdf.Translate3d(v3.Vec{0, 0, -wheelHeight / 2}))
	driver3d = sdf.Union3D(driver3d, base3d)
	// remove a center hole
	hole3d, _ := sdf.Cylinder3D(2*wheelHeight, holeRadius, 0)
	driver3d = sdf.Difference3D(driver3d, hole3d)

	// extrude the driven wheel
	driven3d := sdf.Extrude3D(sDriven, wheelHeight)
	driven3d = sdf.Transform3D(driven3d, sdf.Translate3d(v3.Vec{0, 0, -wheelHeight / 2}))
	// add a hub
	hub3d, _ := sdf.Cylinder3D(wheelHeight, hubRadius, 0)
	hub3d = sdf.Transform3D(hub3d, sdf.Translate3d(v3.Vec{0, 0, wheelHeight / 2}))
	driven3d = sdf.Union3D(driven3d, hub3d)
	// remove a center hole
	driven3d = sdf.Difference3D(driven3d, hole3d)

	const meshCells = 300
	render.ToSTL(driver3d, "driver.stl", render.NewMarchingCubesOctree(meshCells))
	render.ToSTL(driven3d, "driven.stl", render.NewMarchingCubesOctree(meshCells))

	driver3d = sdf.Transform3D(driver3d, sdf.Translate3d(v3.Vec{-0.8 * k.DrivenRadius, 0, 0}))
	driven3d = sdf.Transform3D(driven3d, sdf.Translate3d(v3.Vec{k.DrivenRadius, 0, 0}))
	render.ToSTL(sdf.Union3D(driver3d, driven3d), "geneva.stl", render.NewMarchingCubesOctree(meshCells))
}

//-----------------------------------------------------------------------------
