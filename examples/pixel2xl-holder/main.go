// pixel2xl-holder is an angled holder for a Pixel 2XL that hangs on
// the side wall of a cubicle.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/gmlewis/sdfx/render"
	. "github.com/gmlewis/sdfx/sdf"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

var (
	sampling = flag.Int("sampling", 200, "Number of cells to sample along the longest axis")
	filename = flag.String("filename", "pixel2xl-holder.stl", "Output STL filename")
)

// All dimensions in mm
const (
	mountWall    = 36.5
	phoneWidth   = 82
	phoneDepth   = 11.5
	holderHeight = 104
	// phoneGrab is the amount to overlap the phone in fron in order to hold it in place.
	phoneGrab = 5
	// trumpetStart is the starting height of the trumpet curve where the phone slides in.
	trumpetStart = 20
	// cableDelta is half the cable cutout width.
	cableDelta = 10

	// 3D-printing tollerance
	tol = 0.5
	// 3D-printing wall thickness
	wall = 3
	// edge rounding amount
	round = 2

	// roll is the backward tilt angle (in degrees) of the phone from vertical.
	roll = 10
	// yaw is the backward rotation angle (in degrees) on the right wall of the
	// cubicle from facing forward.
	yaw = 30
)

type V3 = v3.Vec

func main() {
	box, err := Box3D(V3{X: phoneWidth + 2*(tol+wall), Y: phoneDepth + 2*(tol+wall), Z: holderHeight + tol + wall}, round)
	must(err)
	cutout1, err := Box3D(V3{X: phoneWidth + 2*tol, Y: phoneDepth + 2*tol, Z: holderHeight + tol + wall}, round)
	must(err)
	cutout1 = Transform3D(cutout1, Translate3d(V3{X: 0, Y: 0, Z: tol + wall}))
	cutout2, err := Box3D(V3{X: phoneWidth + 2*tol, Y: phoneDepth + 2*(tol+wall), Z: holderHeight + tol + wall}, round)
	must(err)
	cutout2 = Transform3D(cutout2, Translate3d(V3{X: 0, Y: -(tol + wall), Z: tol + wall + phoneGrab}))
	box = Difference3D(box, cutout1)
	box = Difference3D(box, cutout2)

	splinePts := []V2{{X: 0, Y: 0}, {X: 0, Y: trumpetStart}, {X: -5, Y: trumpetStart + 5}, {X: -10, Y: trumpetStart + 10}}
	bez := NewBezierSpline(splinePts)
	poly := NewPolygon()
	bez.Sample(poly, 0, 1, splinePts[0], splinePts[3], 0)
	outline := poly.Vertices()
	for i := len(outline) - 1; i >= 0; i-- {
		p := outline[i]
		outline = append(outline, V2{X: p.X + wall, Y: p.Y})
	}
	polyOutline, err := Polygon2D(outline)
	must(err)
	trumpet := Extrude3D(polyOutline, phoneGrab)
	xfrmL := Translate3d(V3{X: -0.5 * phoneWidth, Y: -0.5*phoneDepth - tol - wall, Z: -0.5 * holderHeight}).Mul(RotateX(0.5 * math.Pi).Mul(RotateY(0.5 * math.Pi)))
	trumpetL := Transform3D(trumpet, xfrmL)
	xfrmR := Translate3d(V3{X: 0.5 * phoneWidth, Y: -0.5*phoneDepth - tol - wall, Z: -0.5 * holderHeight}).Mul(RotateX(0.5 * math.Pi).Mul(RotateY(0.5 * math.Pi)))
	trumpetR := Transform3D(trumpet, xfrmR)
	box = Union3D(box, trumpetL)
	box = Union3D(box, trumpetR)

	wallOutline := poly.Vertices()
	wallOutline = append(wallOutline, V2{X: wall, Y: trumpetStart + 10})
	wallOutline = append(wallOutline, V2{X: wall, Y: 0})
	polyWallOutline, err := Polygon2D(wallOutline)
	must(err)
	wallExtrude := Extrude3D(polyWallOutline, wall)
	xfrmL2 := Translate3d(V3{X: -0.5*(phoneWidth+wall) - tol, Y: -0.5*phoneDepth - tol - wall, Z: -0.5 * holderHeight}).Mul(RotateX(0.5 * math.Pi).Mul(RotateY(0.5 * math.Pi)))
	leftWall := Transform3D(wallExtrude, xfrmL2)
	xfrmR2 := Translate3d(V3{X: 0.5*(phoneWidth+wall) + tol, Y: -0.5*phoneDepth - tol - wall, Z: -0.5 * holderHeight}).Mul(RotateX(0.5 * math.Pi).Mul(RotateY(0.5 * math.Pi)))
	rightWall := Transform3D(wallExtrude, xfrmR2)
	box = Union3D(box, leftWall)
	box = Union3D(box, rightWall)

	bottomCutout := []V2{
		{X: -cableDelta, Y: 0.5*phoneDepth + tol},
		{X: cableDelta, Y: 0.5*phoneDepth + tol},
	}
	r := phoneDepth + 2*tol + wall
	for i := -176; i <= -90; i += 4 {
		bottomCutout = append(bottomCutout, V2{
			X: cableDelta + r + r*math.Cos(float64(i)*math.Pi/180.0),
			Y: 0.5*phoneDepth + tol + r*math.Sin(float64(i)*math.Pi/180.0),
		})
	}
	for i := -90; i >= -176; i -= 4 {
		bottomCutout = append(bottomCutout, V2{
			X: -cableDelta - r - r*math.Cos(float64(i)*math.Pi/180.0),
			Y: 0.5*phoneDepth + tol + r*math.Sin(float64(i)*math.Pi/180.0),
		})
	}
	polyBottomCutout, err := Polygon2D(bottomCutout)
	must(err)
	bottomCutoutExtrude := Extrude3D(polyBottomCutout, 40)
	bottomCutoutExtrude = Transform3D(bottomCutoutExtrude, Translate3d(V3{X: 0, Y: 0, Z: -0.5 * holderHeight}))
	box = Difference3D(box, bottomCutoutExtrude)
	rot := RotateZ(-yaw * math.Pi / 180.0).Mul(RotateX(-roll * math.Pi / 180.0).Mul(Translate3d(V3{X: -phoneWidth - tol - wall, Y: -0.5*phoneDepth - tol - wall, Z: -holderHeight})))
	box = Transform3D(box, rot)

	render.ToSTL(box, *filename, render.NewMarchingCubesOctree(*sampling))
	fmt.Println("Done.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
