package half_utron

import (
	"log"
	"math"

	. "github.com/gmlewis/sdfx/sdf"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

// All dimensions in mm
const (
	minThickness = 3.0
)

type V3 = v3.Vec

func HalfUtron(utronEdge float64) SDF3 {
	cr := math.Sqrt(0.5 * utronEdge * utronEdge)
	cone, err := Cone3D(cr, 0, cr, 0.5)
	must(err)
	cone = Transform3D(cone, Rotate3d(V3{1, 0, 0}, math.Pi))
	cone = Transform3D(cone, Translate3d(V3{0, 0, 0.5 * cr}))

	sd := utronEdge - 2.0*minThickness
	sphere, err := Sphere3D(0.5 * sd)
	must(err)

	return Difference3D(cone, sphere)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
