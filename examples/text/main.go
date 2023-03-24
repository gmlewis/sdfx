//-----------------------------------------------------------------------------
/*

Text Example

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"

	"github.com/gmlewis/sdfx/render"
	"github.com/gmlewis/sdfx/sdf"
	"github.com/gmlewis/sdfx/vec/v2i"
)

//-----------------------------------------------------------------------------

type V2i = v2i.Vec

func main() {

	f, err := sdf.LoadFont("../../files/cmr10.ttf")
	//f, err := sdf.LoadFont("Times_New_Roman.ttf")
	//f, err := sdf.LoadFont("wt064.ttf")

	if err != nil {
		log.Fatalf("can't read font file %s\n", err)
	}

	t := sdf.NewText("SDFX!\nHello,\nWorld!")
	//t := sdf.NewText("相同的不同")

	s2d, err := sdf.TextSDF2(f, t, 10.0)
	if err != nil {
		log.Fatalf("can't generate text sdf2 %s\n", err)
	}

	render.ToDXF(s2d, "shape.dxf", render.NewMarchingSquaresQuadtree(600))
	render.ToSVG(s2d, "shape.svg", render.NewMarchingSquaresQuadtree(600))

	fmt.Println("rendering shape.png (600x525)")
	png, err := render.NewPNG("shape.png", s2d.BoundingBox(), V2i{600, 525})
	if err != nil {
		log.Fatalf("NewPNG: %v", err)
	}
	png.RenderSDF2(s2d)
	if err := png.Save(); err != nil {
		log.Fatalf("Save: %v", err)
	}

	s3d, err := sdf.ExtrudeRounded3D(s2d, 1.0, 0.2)
	if err != nil {
		log.Fatal(err)
	}
	render.ToSTL(s3d, "shape.stl", render.NewMarchingCubesOctree(600))
}

//-----------------------------------------------------------------------------
