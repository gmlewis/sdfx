//-----------------------------------------------------------------------------
/*

Pottery Wheel

*/
//-----------------------------------------------------------------------------

package main

import (
	"log"
	"math"

	"github.com/gmlewis/sdfx/render"
	"github.com/gmlewis/sdfx/sdf"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------
// overall build controls

const shrink = 1.0 / 0.98 // 2% Al shrinkage
const core_print = false  // add the core print to the wheel
const pie_print = false   // create a 1/n pie segment (n = number of webs)

//-----------------------------------------------------------------------------

// draft angles
var draft_angle = sdf.DtoR(4.0)       // standard overall draft
var core_draft_angle = sdf.DtoR(10.0) // draft angle for the core print

// nominal size values (mm)
const wheel_diameter = sdf.MillimetresPerInch * 8.0 // total wheel diameter
const hub_diameter = 40.0                           // base diameter of central shaft hub
const hub_height = 53.0                             // height of cental shaft hub
const shaft_diameter = 21.0                         // 1" target size - reduced for machining allowance
const shaft_length = 45.0                           // length of shaft bore
const wall_height = 35.0                            // height of wheel side walls
const wall_thickness = 4.0                          // base thickness of outer wheel walls
const plate_thickness = 7.0                         // thickness of wheel top plate
const web_width = 2.0                               // thickness of reinforcing webs
const web_height = 25.0                             // height of reinforcing webs
const core_height = 15.0                            // height of core print
const number_of_webs = 6                            // number of reinforcing webs

// derived values
const wheel_radius = wheel_diameter / 2
const hub_radius = hub_diameter / 2
const shaft_radius = shaft_diameter / 2
const web_length = wheel_radius - wall_thickness - hub_radius

//-----------------------------------------------------------------------------

// build 2d wheel profile
func wheel_profile() (sdf.SDF2, error) {

	draft0 := (hub_height - plate_thickness) * math.Tan(draft_angle)
	draft1 := (wall_height - plate_thickness) * math.Tan(draft_angle)
	draft2 := wall_height * math.Tan(draft_angle)
	draft3 := core_height * math.Tan(core_draft_angle)

	s := sdf.NewPolygon()

	if core_print {
		s.Add(0, 0)
		s.Add(wheel_radius+draft2, 0)
		s.Add(wheel_radius, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness-draft1, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius+draft0, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius, hub_height).Smooth(2.0, 5)
		s.Add(shaft_radius, hub_height)
		s.Add(shaft_radius-draft3, hub_height+core_height)
		s.Add(0, hub_height+core_height)
	} else {
		s.Add(0, 0)
		s.Add(wheel_radius+draft2, 0)
		s.Add(wheel_radius, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness-draft1, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius+draft0, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius, hub_height).Smooth(2.0, 5)
		s.Add(shaft_radius, hub_height)
		s.Add(shaft_radius, hub_height-shaft_length)
		s.Add(0, hub_height-shaft_length)
	}

	return sdf.Polygon2D(s.Vertices())
}

// build 2d web profile
func web_profile() (sdf.SDF2, error) {

	draft := web_height * math.Tan(draft_angle)
	x0 := web_width + draft
	x1 := web_width

	s := sdf.NewPolygon()
	s.Add(-x0, 0)
	s.Add(-x1, web_height).Smooth(1.0, 3)
	s.Add(x1, web_height).Smooth(1.0, 3)
	s.Add(x0, 0)

	//s.Render("web.dxf")
	return sdf.Polygon2D(s.Vertices())
}

// build the wheel pattern
func wheel_pattern() (sdf.SDF3, error) {

	// build a reinforcing webs
	web_2d, err := web_profile()
	if err != nil {
		return nil, err
	}
	web_3d := sdf.Extrude3D(web_2d, web_length)
	m := sdf.Translate3d(v3.Vec{0, plate_thickness, hub_radius + web_length/2})
	m = sdf.RotateX(sdf.DtoR(90)).Mul(m)

	// build the wheel profile
	wheel_2d, err := wheel_profile()
	if err != nil {
		return nil, err
	}

	var wheel_3d sdf.SDF3
	if pie_print {
		m = sdf.RotateZ(sdf.DtoR(120)).Mul(m)
		web_3d = sdf.Transform3D(web_3d, m)
		wheel_3d, err = sdf.RevolveTheta3D(wheel_2d, sdf.DtoR(60))
	} else {
		m = sdf.RotateZ(sdf.DtoR(90)).Mul(m)
		web_3d = sdf.Transform3D(web_3d, m)
		web_3d = sdf.RotateCopy3D(web_3d, 6)
		wheel_3d, err = sdf.Revolve3D(wheel_2d)
	}
	if err != nil {
		return nil, err
	}

	// add the webs to the wheel with some blending
	wheel := sdf.Union3D(wheel_3d, web_3d)
	wheel.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(wall_thickness))
	return wheel, nil
}

//-----------------------------------------------------------------------------

// build 2d core profile
func core_profile() (sdf.SDF2, error) {

	draft := core_height * math.Tan(core_draft_angle)

	s := sdf.NewPolygon()
	s.Add(0, 0)
	s.Add(shaft_radius-draft, 0)
	s.Add(shaft_radius, core_height)
	s.Add(shaft_radius, core_height+shaft_length).Smooth(2.0, 3)
	s.Add(0, core_height+shaft_length)

	//s.Render("core.dxf")
	return sdf.Polygon2D(s.Vertices())
}

// build the core box
func core_box() (sdf.SDF3, error) {

	// build the box
	w := 4.2 * shaft_radius
	d := 1.2 * shaft_radius
	h := (core_height + shaft_length) * 1.1
	box_3d, err := sdf.Box3D(v3.Vec{h, w, d}, 0)
	if err != nil {
		return nil, err
	}

	// holes in the box
	dy := w * 0.37
	dx := h * 0.4
	hole_radius := ((3.0 / 16.0) * sdf.MillimetresPerInch) / 2.0
	positions := []v3.Vec{
		{dx, dy, 0},
		{-dx, dy, 0},
		{dx, -dy, 0},
		{-dx, -dy, 0},
	}
	hole, err := sdf.Cylinder3D(d, hole_radius, 0)
	if err != nil {
		return nil, err
	}
	holes_3d := sdf.Multi3D(hole, positions)

	// Drill the holes
	box_3d = sdf.Difference3D(box_3d, holes_3d)

	// build the core
	core_2d, err := core_profile()
	if err != nil {
		return nil, err
	}

	core_3d, err := sdf.Revolve3D(core_2d)
	if err != nil {
		return nil, err
	}

	m := sdf.Translate3d(v3.Vec{h / 2, 0, d / 2}).Mul(sdf.RotateY(sdf.DtoR(-90)))
	core_3d = sdf.Transform3D(core_3d, m)

	// remove the core from the box
	return sdf.Difference3D(box_3d, core_3d), nil
}

//-----------------------------------------------------------------------------

func main() {
	s0, err := wheel_pattern()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s0 = sdf.ScaleUniform3D(s0, shrink)
	render.ToSTL(s0, "wheel.stl", render.NewMarchingCubesOctree(200))
	render.ToDXF(sdf.Slice2D(s0, v3.Vec{0, 0, 15.0}, v3.Vec{0, 0, 1}), "wheel.dxf", render.NewMarchingSquaresQuadtree(200))

	s1, err := core_box()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s1 = sdf.ScaleUniform3D(s1, shrink)
	render.ToSTL(s1, "core_box.stl", render.NewMarchingCubesOctree(200))
}

//-----------------------------------------------------------------------------
