package enclosure

import (
	"log"
	"math"

	"github.com/gmlewis/sdfx/obj"
	. "github.com/gmlewis/sdfx/sdf"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

// All dimensions in mm
const (
	utronMargin   = 5.0
	dielectricGap = 1.5

	magnetMargin = 10.0
	gapWidth     = 50.0
	innerGap     = 70.0
	magnetHeight = 101.6
	magnetDiam   = 50.8

	baseHeight    = 11.0
	wallThickness = 16.0

	bearingHeight     = 5.0
	bearingDiam       = 14.0
	bearingMarginDiam = 0.75
	bearingMarginZ    = 0.5
	bearingOverhang   = 2.0

	boltDiam   = 0.75 * wallThickness
	boltHeight = 10.0
)

type V3 = v3.Vec

func Top(utronEdge float64) SDF3 {
	utronDiam := math.Sqrt(2 * utronEdge * utronEdge)

	inside := utronDiam + 2*utronMargin
	outside := inside + 2*wallThickness
	boxHeight := baseHeight + 0.5*bearingHeight
	box, err := Box3D(V3{outside, outside, boxHeight}, 0)
	must(err)
	box = Transform3D(box, Translate3d(V3{0, 0, utronDiam + 0.5*boxHeight - 0.5*bearingHeight}))

	dx := 3.0
	wallHeight := 0.5*utronDiam + baseHeight - 0.5*bearingHeight
	walls, err := Box3D(V3{0.5*outside - dx, outside, wallHeight}, 0)
	must(err)
	walls = Transform3D(walls, Translate3d(V3{-0.25*outside - 0.5*dx, 0, utronDiam - 0.5*wallHeight}))
	box = Union3D(box, walls)
	wallCutout, err := Box3D(V3{inside, inside, wallHeight}, 0)
	must(err)
	wallCutout = Transform3D(wallCutout, Translate3d(V3{0, 0, utronDiam - 0.5*wallHeight - 0.5*bearingHeight}))
	box = Difference3D(box, wallCutout)

	// upper magnet brace
	dx = math.Sqrt(2 * utronMargin * utronMargin)
	ts := 0.5*outside - wallThickness
	triangle, err := Polygon2D([]V2{{-dx, utronDiam}, {-ts - dx, utronDiam}, {-ts - dx, utronDiam - ts}})
	must(err)
	prism := Extrude3D(triangle, outside)
	prism = Transform3D(prism, RotateX(0.5*math.Pi))
	box = Union3D(box, prism)

	big := 10 * utronDiam
	cyl, err := Cylinder3D(big, 0.5*(boltDiam+1), 0)
	must(err)
	cyl1 := Transform3D(cyl, Translate3d(hole1(outside, utronDiam)))
	box = Difference3D(box, cyl1)
	cyl2 := Transform3D(cyl, Translate3d(hole2(outside, utronDiam)))
	box = Difference3D(box, cyl2)
	cyl3 := Transform3D(cyl, Translate3d(hole3(outside, utronDiam)))
	box = Difference3D(box, cyl3)
	cyl4 := Transform3D(cyl, Translate3d(hole4(outside, utronDiam)))
	box = Difference3D(box, cyl4)
	cyl5 := Transform3D(cyl, Translate3d(hole5(outside, utronDiam)))
	box = Difference3D(box, cyl5)
	cyl6 := Transform3D(cyl, Translate3d(hole6(outside, utronDiam)))
	box = Difference3D(box, cyl6)

	topDuct, err := Cylinder3D(big, 0.5*0.45*utronDiam, 0)
	must(err)
	x := 0.5 * (0.5*wallThickness + 0.5*(outside-wallThickness))
	topDuct1 := Transform3D(topDuct, Translate3d(V3{x, -0.25 * (outside - wallThickness), utronDiam}))
	topDuct2 := Transform3D(topDuct, Translate3d(V3{x, 0.25 * (outside - wallThickness), utronDiam}))
	box = Difference3D(box, topDuct1)
	box = Difference3D(box, topDuct2)

	bearing, err := Cylinder3D(bearingHeight+2*bearingMarginZ, 0.5*(bearingDiam+bearingMarginDiam), 0)
	must(err)
	bearing = Transform3D(bearing, Translate3d(V3{0, 0, utronDiam}))
	access, err := Cylinder3D(wallThickness, 0.5*(bearingDiam-bearingOverhang), 0)
	must(err)
	access = Transform3D(access, Translate3d(V3{0, 0, utronDiam + 0.5*wallThickness}))
	bearingCutout := Union3D(bearing, access)

	return Difference3D(box, bearingCutout)
}

func hole1(outside, z float64) V3 {
	return V3{0.5 * wallThickness, -0.5 * (outside - wallThickness), z}
}

func hole2(outside, z float64) V3 {
	return V3{0.5 * (outside - wallThickness), -0.5 * (outside - wallThickness), z}
}

func hole3(outside, z float64) V3 {
	return V3{0.5 * wallThickness, 0.5 * (outside - wallThickness), z}
}

func hole4(outside, z float64) V3 {
	return V3{0.5 * (outside - wallThickness), 0.5 * (outside - wallThickness), z}
}

func hole5(outside, z float64) V3 {
	return V3{0.5 * (outside - wallThickness), 0, z}
}

func hole6(outside, z float64) V3 {
	return V3{-0.5 * (outside - wallThickness), -0.5 * (outside - wallThickness), z}
}

func Base(utronEdge float64) SDF3 {
	utronDiam := math.Sqrt(2 * utronEdge * utronEdge)

	// center of lower bearing is the origin.
	inside := utronDiam + 2*utronMargin
	outside := inside + 2*wallThickness
	inbox, err := Box3D(V3{inside, inside, 2 * outside}, 0)
	must(err)
	inbox = Transform3D(inbox, Translate3d(V3{0, 0, outside}))
	boxHeight := wallThickness - 1.5*bearingHeight + utronDiam
	box, err := Box3D(V3{outside, outside, boxHeight}, 0)
	must(err)
	box = Transform3D(box, Translate3d(V3{0, 0, 0.5*boxHeight - wallThickness}))
	box = Difference3D(box, inbox)
	box = Transform3D(box, Translate3d(V3{0, 0, 0.5 * bearingHeight}))
	// left cutout
	cutBox, err := Box3D(V3{outside, outside, outside}, 0)
	must(err)
	cutPosZ := 0.5*utronDiam - baseHeight
	cutBox = Transform3D(cutBox, Translate3d(V3{-0.5 * outside, 0, 0.5*outside + cutPosZ}))
	box = Difference3D(box, cutBox)
	// lower magnet brace
	dx := math.Sqrt(2 * utronMargin * utronMargin)
	ts := 0.5*outside - wallThickness
	triangle, err := Polygon2D([]V2{{dx, 0}, {ts + dx, 0}, {ts + dx, ts}})
	must(err)
	prism := Extrude3D(triangle, outside)
	prism = Transform3D(prism, RotateX(0.5*math.Pi))
	box = Union3D(box, prism)

	boxTopZ := utronDiam - bearingHeight
	h := baseHeight + bearingHeight
	box = addBolt(box, h, hole1(outside, boxTopZ))
	box = addBolt(box, h, hole2(outside, boxTopZ))
	box = addBolt(box, h, hole3(outside, boxTopZ))
	box = addBolt(box, h, hole4(outside, boxTopZ))
	box = addBolt(box, h, hole5(outside, boxTopZ))
	h = 0.5*utronDiam + 2*baseHeight
	box = addBolt(box, h, hole6(outside, cutPosZ))

	// air ducts.
	airDuct, err := Cylinder3D(outside, utronDiam/6, 0)
	must(err)
	airDuct = Transform3D(airDuct, RotateX(0.5*math.Pi))
	airDuct = Transform3D(airDuct, Translate3d(V3{0.25 * outside, 0, boxTopZ - utronDiam/3}))
	box = Difference3D(box, airDuct)
	sideDuct, err := Cylinder3D(outside, utronDiam/6, 0)
	must(err)
	sideDuct = Transform3D(sideDuct, RotateY(0.5*math.Pi))
	sideDuct1 := Transform3D(sideDuct, Translate3d(V3{0, -0.25 * (outside - wallThickness), boxTopZ - utronDiam/3}))
	sideDuct2 := Transform3D(sideDuct, Translate3d(V3{0, 0.25 * (outside - wallThickness), boxTopZ - utronDiam/3}))
	box = Difference3D(box, sideDuct1)
	box = Difference3D(box, sideDuct2)

	bearing, err := Cylinder3D(bearingHeight+2*bearingMarginZ, 0.5*(bearingDiam+bearingMarginDiam), 0)
	must(err)
	access, err := Cylinder3D(wallThickness, 0.5*(bearingDiam-bearingOverhang), 0)
	must(err)
	access = Transform3D(access, Translate3d(V3{0, 0, -0.5 * wallThickness}))
	bearingCutout := Union3D(bearing, access)

	return Difference3D(box, bearingCutout)
}

func addBolt(box SDF3, height float64, basePos V3) SDF3 {
	// shaft := Cylinder3D(height, 0.5*boltDiam, 0)
	// shaft = Transform3D(shaft, Translate3d(basePos.Add(V3{0, 0, 0.5 * height})))
	h := dielectricGap + 1.5*boltHeight
	// threads := Cylinder3D(h, 0.5*boltDiam, 0)
	// threads = Transform3D(threads, Translate3d(basePos.Add(V3{0, 0, 0.5*h + height})))
	// bolt := Union3D(shaft, threads)

	overlap := 1.0 // remove chamfer at connection point
	bolt := Nut_And_Bolt("M12x1.5", 0, overlap+height+h, overlap+height)
	bolt = Transform3D(bolt, Translate3d(basePos.Add(V3{0, 0, -overlap})))

	return Union3D(box, bolt)
}

//////////////////////////////////////////////////////////////////////
// Nuts and bolts taken from nutsandbolts example
//////////////////////////////////////////////////////////////////////

func Hex_Bolt(
	name string, // name of thread
	tolerance float64, // subtract from external thread radius
	total_length float64, // threaded length + shank length
	shank_length float64, //  non threaded length
) SDF3 {

	t, err := ThreadLookup(name)
	must(err)

	if total_length < 0 {
		return nil
	}
	if shank_length < 0 {
		return nil
	}
	thread_length := total_length - shank_length
	if thread_length < 0 {
		thread_length = 0
	}

	// 	// hex head
	hex_r := t.HexRadius()
	// hex_h := t.Hex_Height()
	// hex_3d := HexHead3D(hex_r, hex_h, "b")
	//
	// 	// add a rounded cylinder
	// 	hex_3d = Union3D(hex_3d, Cylinder3D(hex_h*1.05, hex_r*0.8, hex_r*0.08))

	// shank
	// shank_length += hex_h / 2
	shank_ofs := shank_length / 2
	shank_3d, err := Cylinder3D(shank_length, t.Radius, hex_r*0.08)
	must(err)
	shank_3d = Transform3D(shank_3d, Translate3d(V3{0, 0, shank_ofs}))

	// thread
	r := t.Radius - tolerance
	l := thread_length
	screw_ofs := l/2 + shank_length
	extIso, err := ISOThread(r, t.Pitch, true)
	must(err)
	screw_3d, err := Screw3D(extIso, l, 0, t.Pitch, 1)
	must(err)
	// chamfer the thread
	screw_3d, err = obj.ChamferedCylinder(screw_3d, 0, 0.5)
	must(err)
	screw_3d = Transform3D(screw_3d, Translate3d(V3{0, 0, screw_ofs}))

	// return Union3D(hex_3d, screw_3d, shank_3d)
	return Union3D(screw_3d, shank_3d)
}

//-----------------------------------------------------------------------------

// Return a Hex Nut
func Hex_Nut(
	name string, // name of thread
	tolerance float64, // add to internal thread radius
	height float64, // height of nut
) SDF3 {

	t, err := ThreadLookup(name)
	must(err)

	if height < 0 {
		return nil
	}

	// hex nut body
	hex_3d, err := obj.HexHead3D(t.HexRadius(), height, "tb")
	must(err)

	// internal thread
	intIso, err := ISOThread(t.Radius+tolerance, t.Pitch, false)
	must(err)
	thread_3d, err := Screw3D(intIso, height, 0, t.Pitch, 1)
	must(err)

	return Difference3D(hex_3d, thread_3d)
}

//-----------------------------------------------------------------------------

func Nut_And_Bolt(
	name string, // name of thread
	tolerance float64, // thread tolerance
	total_length float64, // threaded length + shank length
	shank_length float64, //  non threaded length
) SDF3 {
	t, err := ThreadLookup(name)
	must(err)
	bolt_3d := Hex_Bolt(name, tolerance, total_length, shank_length)
	nut_3d := Hex_Nut(name, tolerance, t.HexHeight()/1.5)
	z_ofs := total_length + t.HexHeight() + 0.25
	nut_3d = Transform3D(nut_3d, Translate3d(V3{0, 0, z_ofs}))
	return Union3D(nut_3d, bolt_3d)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
