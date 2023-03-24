// -*- compile-command: "go build && ./utron && fstl utron.stl"; -*-

package main

import (
	"log"
	"math"

	"github.com/gmlewis/sdfx/examples/utron/enclosure"
	half_magnet "github.com/gmlewis/sdfx/examples/utron/half-magnet"
	half_utron "github.com/gmlewis/sdfx/examples/utron/half-utron"
	"github.com/gmlewis/sdfx/render"
	. "github.com/gmlewis/sdfx/sdf"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

// All dimensions in mm
const (
	utronEdge    = 50.0
	magnetHeight = 25.4
	innerGap     = 70.0
	magnetDiam   = 50.8
	metalMargin  = 0.5
	magnetMargin = 10.0
)

var (
	utronRadius = 0.5 * math.Sqrt(2*utronEdge*utronEdge)
)

type V3 = v3.Vec

func top() SDF3 {
	top := enclosure.Top(utronEdge)
	ch := 4 * magnetHeight
	topCutout, err := Cylinder3D(ch, 0.5*magnetDiam+metalMargin, 1)
	must(err)
	ssHeight := 0.5*(4*magnetHeight-utronEdge) - magnetMargin
	m := Translate3d(V3{0, 0, 0.5*ch + 2*magnetHeight - ssHeight - metalMargin})
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	topCutout = Transform3D(topCutout, m)
	side := magnetDiam + 2*metalMargin
	big := 10 * utronEdge
	boxCutout, err := Box3D(V3{side, big, side}, 0)
	must(err)
	m = Translate3d(V3{0, 0.5 * big, 0.5*side + 2*magnetHeight - ssHeight - metalMargin})
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	boxCutout = Transform3D(boxCutout, m)
	topCutout = Union3D(topCutout, boxCutout)
	top = Difference3D(top, topCutout)

	return top
}

func base() SDF3 {
	base := enclosure.Base(utronEdge)
	ch := 4 * magnetHeight
	baseCutout, err := Cylinder3D(ch, 0.5*magnetDiam+metalMargin, 1)
	must(err)
	ssHeight := 0.5*(4*magnetHeight-utronEdge) - magnetMargin
	m := Translate3d(V3{0, 0, -0.5*ch - 2*magnetHeight + ssHeight + metalMargin})
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	baseCutout = Transform3D(baseCutout, m)
	side := magnetDiam + 2*metalMargin
	big := 10 * utronEdge
	boxCutout, err := Box3D(V3{side, big, side}, 0)
	must(err)
	m = Translate3d(V3{0, 0.5 * big, -0.5*side - 2*magnetHeight + ssHeight + metalMargin})
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	boxCutout = Transform3D(boxCutout, m)
	baseCutout = Union3D(baseCutout, boxCutout)
	base = Difference3D(base, baseCutout)

	return base
}

func main() {
	top := top()
	base := base()

	halfUtron := half_utron.HalfUtron(utronEdge)
	utronLower := Transform3D(halfUtron, RotateX(math.Pi))
	utronLower = Transform3D(utronLower, Translate3d(V3{0, 0, utronRadius}))
	utronUpper := Transform3D(halfUtron, Translate3d(V3{0, 0, utronRadius}))

	halfMagnet := half_magnet.HalfMagnet(utronEdge, innerGap, magnetDiam, magnetHeight, magnetMargin)
	m := RotateX(0.5 * math.Pi)
	m = Translate3d(V3{-0.5 * (innerGap + magnetDiam), 0, -2 * magnetHeight}).Mul(m)
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	halfMagnetLower := Transform3D(halfMagnet, m)
	m = RotateX(-0.5 * math.Pi)
	m = Translate3d(V3{-0.5 * (innerGap + magnetDiam), 0, 2 * magnetHeight}).Mul(m)
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	halfMagnetUpper := Transform3D(halfMagnet, m)

	trim := 1.0 // To separate each magnet into its own piece and prevent merging.
	magnet1, err := Cylinder3D(magnetHeight-trim, 0.5*magnetDiam, 1)
	must(err)
	magnet1 = Transform3D(magnet1, Translate3d(V3{0, 0, -1.5 * magnetHeight}))
	magnet2, err := Cylinder3D(magnetHeight-trim, 0.5*magnetDiam, 1)
	must(err)
	magnet2 = Transform3D(magnet2, Translate3d(V3{0, 0, -0.5 * magnetHeight}))
	magnet3, err := Cylinder3D(magnetHeight-trim, 0.5*magnetDiam, 1)
	must(err)
	magnet3 = Transform3D(magnet3, Translate3d(V3{0, 0, 0.5 * magnetHeight}))
	magnet4, err := Cylinder3D(magnetHeight-trim, 0.5*magnetDiam, 1)
	must(err)
	magnet4 = Transform3D(magnet4, Translate3d(V3{0, 0, 1.5 * magnetHeight}))
	magnets := Union3D(magnet1, magnet2, magnet3, magnet4)
	m = Translate3d(V3{-innerGap - magnetDiam, 0, 0})
	m = RotateY(-0.25 * math.Pi).Mul(m)
	m = Translate3d(V3{0, 0, utronRadius}).Mul(m)
	magnets = Transform3D(magnets, m)

	s := Union3D(base, utronLower, utronUpper, halfMagnetLower, halfMagnetUpper, magnets, top)
	render.ToSTL(s, "utron.stl", render.NewMarchingCubesOctree(800))

	// Write out separate parts.
	render.ToSTL(base, "base.stl", render.NewMarchingCubesOctree(800))
	render.ToSTL(top, "top.stl", render.NewMarchingCubesOctree(800))
	render.ToSTL(utronLower, "utron-lower.stl", render.NewMarchingCubesOctree(800))
	render.ToSTL(utronUpper, "utron-upper.stl", render.NewMarchingCubesOctree(800))
	render.ToSTL(halfMagnetLower, "magnet-lower.stl", render.NewMarchingCubesOctree(800))
	render.ToSTL(halfMagnetUpper, "magnet-upper.stl", render.NewMarchingCubesOctree(800))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
