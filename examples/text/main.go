//-----------------------------------------------------------------------------
/*

Text Example

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"log"
	"os"

	. "github.com/gmlewis/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {

	f, err := LoadFont("cmr10.ttf")
	//f, err := LoadFont("Times_New_Roman.ttf")
	//f, err := LoadFont("wt064.ttf")

	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}

	t := NewText("SDFX!\nHello,\nWorld!")
	//t := NewText("相同的不同")

	s2d, err := TextSDF2(f, t, 10.0)
	if err != nil {
		fmt.Printf("can't generate text sdf2 %s\n", err)
		os.Exit(1)
	}

	RenderDXF(s2d, 600, "shape.dxf")
	RenderSVG(s2d, 600, "shape.svg", "fill:none;stroke:black;stroke-width:0.1")

	fmt.Println("rendering shape.png (600x525)")
	png, err := NewPNG("shape.png", s2d.BoundingBox(), V2i{600, 525})
	if err != nil {
		log.Fatalf("NewPNG: %v", err)
	}
	png.RenderSDF2(s2d)
	if err := png.Save(); err != nil {
		log.Fatalf("Save: %v", err)
	}

	s3d := ExtrudeRounded3D(s2d, 1.0, 0.2)
	RenderSTL(s3d, 600, "shape.stl")
}

//-----------------------------------------------------------------------------
