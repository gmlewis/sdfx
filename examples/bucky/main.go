//-----------------------------------------------------------------------------
/*

Bucky Ball

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

type edge [2]int

var edges = []edge{
	{1, 4}, {4, 8}, {8, 10}, {10, 6}, {6, 1},
	{0, 4}, {4, 9}, {9, 11}, {11, 6}, {6, 0},
	{3, 7}, {7, 10}, {10, 8}, {8, 5}, {5, 3},
	{5, 9}, {9, 11}, {11, 7}, {7, 2}, {2, 5},
	{0, 1}, {1, 9}, {9, 5}, {5, 8}, {8, 0},
	{4, 9}, {9, 3}, {3, 2}, {2, 8}, {8, 4},
	{1, 0}, {0, 10}, {10, 7}, {7, 11}, {11, 1},
	{2, 3}, {3, 11}, {11, 6}, {6, 10}, {10, 2},
	{0, 10}, {10, 2}, {2, 5}, {5, 4}, {4, 0},
	{1, 4}, {4, 5}, {5, 3}, {3, 11}, {11, 1},
	{7, 2}, {2, 8}, {8, 0}, {0, 6}, {6, 7},
	{1, 6}, {6, 7}, {7, 3}, {3, 9}, {9, 1},
}

const φ = 1.618033988749895

var vertex = []v3.Vec{
	{1, φ, 0}, {-1, φ, 0}, {1, -φ, 0}, {-1, -φ, 0},
	{0, 1, φ}, {0, -1, φ}, {0, 1, -φ}, {0, -1, -φ},
	{φ, 0, 1}, {-φ, 0, 1}, {φ, 0, -1}, {-φ, 0, -1},
}

func icosahedron() (sdf.SDF3, error) {

	r0 := φ * 0.05
	r1 := r0 * 2.0

	k := obj.ArrowParms{
		Axis:  [2]float64{0, r0},
		Head:  [2]float64{0, r1},
		Tail:  [2]float64{0, r1},
		Style: "b.",
	}

	var bb sdf.SDF3
	for _, e := range edges {
		head := vertex[e[0]]
		tail := vertex[e[1]]
		s, err := obj.DirectedArrow3D(&k, head, tail)
		if err != nil {
			return nil, err
		}
		bb = sdf.Union3D(bb, s)
	}

	return bb, nil
}

//-----------------------------------------------------------------------------

func main() {
	s, err := icosahedron()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	render.ToSTL(s, "icosahedron.stl", render.NewMarchingCubesOctree(300))
	render.To3MF(s, "icosahedron.3mf", render.NewMarchingCubesOctree(300))
}

//-----------------------------------------------------------------------------
