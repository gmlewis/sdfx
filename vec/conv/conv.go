//-----------------------------------------------------------------------------
/*

Vector Conversions

*/
//-----------------------------------------------------------------------------

package conv

import (
	"math"

	"github.com/gmlewis/sdfx/vec/p2"
	v2 "github.com/gmlewis/sdfx/vec/v2"
	"github.com/gmlewis/sdfx/vec/v2i"
	v3 "github.com/gmlewis/sdfx/vec/v3"
	"github.com/gmlewis/sdfx/vec/v3i"
)

//-----------------------------------------------------------------------------
// V2i to X

// V2iToV2 converts a 2D integer vector to a float vector.
func V2iToV2(a v2i.Vec) v2.Vec {
	return v2.Vec{float64(a.X), float64(a.Y)}
}

//-----------------------------------------------------------------------------
// V3i to X

// V3iToV3 converts a 3D integer vector to a float vector.
func V3iToV3(a v3i.Vec) v3.Vec {
	return v3.Vec{float64(a.X), float64(a.Y), float64(a.Z)}
}

//-----------------------------------------------------------------------------
// V2 to X

// V2ToP2 converts a cartesian to a polar coordinate.
func V2ToP2(a v2.Vec) p2.Vec {
	return p2.Vec{a.Length(), math.Atan2(a.Y, a.X)}
}

// V2ToV3 converts a 2D vector to a 3D vector with a specified Z value.
func V2ToV3(a v2.Vec, z float64) v3.Vec {
	return v3.Vec{a.X, a.Y, z}
}

// V2ToV2i converts a 2D float vector to a 2D integer vector.
func V2ToV2i(a v2.Vec) v2i.Vec {
	return v2i.Vec{int(a.X), int(a.Y)}
}

//-----------------------------------------------------------------------------
// V3 to X

// V3ToV3i converts a 3D float vector to a 3D integer vector.
func V3ToV3i(a v3.Vec) v3i.Vec {
	return v3i.Vec{int(a.X), int(a.Y), int(a.Z)}
}

//-----------------------------------------------------------------------------
// P2 to X

// P2ToV2 converts a polar to a cartesian coordinate.
func P2ToV2(a p2.Vec) v2.Vec {
	return v2.Vec{a.R * math.Cos(a.Theta), a.R * math.Sin(a.Theta)}
}

//-----------------------------------------------------------------------------
