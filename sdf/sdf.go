//-----------------------------------------------------------------------------
/*

 */
//-----------------------------------------------------------------------------

package sdf

import (
	"math"

	"github.com/deadsy/pt/pt"
)

//-----------------------------------------------------------------------------

type SDF3 interface {
	Evaluate(p V3) float64
	BoundingBox() Box3
}

type SDF2 interface {
	Evaluate(p V2) float64
	BoundingBox() Box2
}

//-----------------------------------------------------------------------------
// Basic SDF Functions

func sdf_box3d(p, s V3) float64 {
	d := p.Abs().Sub(s)
	return d.Max(V3{0, 0, 0}).Length() + math.Min(d.MaxComponent(), 0)
}

func sdf_box2d(p, s V2) float64 {
	d := p.Abs().Sub(s)
	return d.Max(V2{0, 0}).Length() + math.Min(d.MaxComponent(), 0)
}

/* alternate function - probably faster
func sdf_box2d(p, s V2) float64 {
	p = p.Abs()
	d := p.Sub(s)
	k := s.Y - s.X
	if d.X > 0 && d.Y > 0 {
		return d.Length()
	}
	if p.Y-p.X > k {
		return d.Y
	}
	return d.X
}
*/

//-----------------------------------------------------------------------------
// Minimum Functions

type MinFunc func(a, b, k float64) float64

// normal min - no blending
func NormalMin(a, b, k float64) float64 {
	return math.Min(a, b)
}

// round min uses a quarter-circle to join the two objects smoothly
func RoundMin(a, b, k float64) float64 {
	u := V2{k - a, k - b}.Max(V2{0, 0})
	return math.Max(k, math.Min(a, b)) - u.Length()
}

// chamfer min makes a 45-degree chamfered edge (the diagonal of a square of size <r>)
func ChamferMin(a, b, k float64) float64 {
	return math.Min(math.Min(a, b), (a-k+b)*math.Sqrt(0.5))
}

// exponential smooth min (k = 32);
func ExpMin(a, b, k float64) float64 {
	return -math.Log(math.Exp(-k*a)+math.Exp(-k*b)) / k
}

// power smooth min (k = 8) (TODO - weird results, is this correct?)
func PowMin(a, b, k float64) float64 {
	a = math.Pow(a, k)
	b = math.Pow(b, k)
	return math.Pow((a*b)/(a+b), 1/k)
}

// polynomial smooth min (k = 0.1)
func PolyMin(a, b, k float64) float64 {
	h := Clamp(0.5+0.5*(b-a)/k, 0.0, 1.0)
	return Mix(b, a, h) - k*h*(1.0-h)
}

//-----------------------------------------------------------------------------
// Create a pt.SDF from an SDF3

type PtSDF struct {
	Sdf SDF3
}

func NewPtSDF(sdf SDF3) pt.SDF {
	return &PtSDF{sdf}
}

func (s *PtSDF) Evaluate(p pt.Vector) float64 {
	return s.Sdf.Evaluate(V3{p.X, p.Y, p.Z})
}

func (s *PtSDF) BoundingBox() pt.Box {
	b := s.Sdf.BoundingBox()
	j := b.Min
	k := b.Max
	return pt.Box{Min: pt.Vector{X: j.X, Y: j.Y, Z: j.Z}, Max: pt.Vector{X: k.X, Y: k.Y, Z: k.Z}}
}

//-----------------------------------------------------------------------------
// Solid of Revolution, SDF2 -> SDF3

type SorSDF3 struct {
	Sdf   SDF2
	Theta float64 // angle for partial revolutions
	Norm  V2      // pre-calculated normal to theta line
}

func NewSorSDF3(sdf SDF2) SDF3 {
	return &SorSDF3{sdf, 0, V2{}}
}

func NewSorThetaSDF3(sdf SDF2, theta float64) SDF3 {
	// normalize theta
	theta = math.Mod(math.Abs(theta), TAU)
	// pre-calculate the normal to the theta line
	norm := V2{math.Sin(theta), -math.Cos(theta)}
	return &SorSDF3{sdf, theta, norm}
}

func (s *SorSDF3) Evaluate(p V3) float64 {
	x := math.Sqrt(p.X*p.X + p.Y*p.Y)
	a := s.Sdf.Evaluate(V2{x, p.Z})
	b := a
	if s.Theta != 0 {
		// combine two vertical planes to give an intersection wedge
		d := s.Norm.Dot(V2{p.X, p.Y})
		if s.Theta < PI {
			b = math.Max(p.Y, d) // intersect
		} else {
			b = math.Min(p.Y, d) // union
		}
	}
	// return the intersection
	return math.Max(a, b)
}

func (s *SorSDF3) BoundingBox() Box3 {
	// TODO - reduce the BB for theta != 0
	b := s.Sdf.BoundingBox()
	j := b.Min
	k := b.Max
	l := math.Max(math.Abs(j.X), math.Abs(k.X))
	return Box3{V3{-l, -l, j.Y}, V3{l, l, k.Y}}
}

//-----------------------------------------------------------------------------
// Extrude, SDF2 -> SDF3

type ExtrudeSDF3 struct {
	Sdf    SDF2
	Height float64
}

func NewExtrudeSDF3(sdf SDF2, height float64) SDF3 {
	return &ExtrudeSDF3{sdf, height}
}

func (s *ExtrudeSDF3) Evaluate(p V3) float64 {
	// sdf for the projected 2d surface
	a := s.Sdf.Evaluate(V2{p.X, p.Y})
	// sdf for the extrusion region: z = [0, height]
	b := math.Max(-p.Z, p.Z-s.Height)
	// return the intersection
	return math.Max(a, b)
}

func (s *ExtrudeSDF3) BoundingBox() Box3 {
	b := s.Sdf.BoundingBox()
	j := b.Min
	k := b.Max
	return Box3{V3{j.X, j.Y, 0}, V3{k.X, k.Y, s.Height}}
}

//-----------------------------------------------------------------------------
// 3D Normal Box

type BoxSDF3 struct {
	Size V3
}

func NewBoxSDF3(size V3) SDF3 {
	// note: store a modified size
	return &BoxSDF3{size.MulScalar(0.5)}
}

func (s *BoxSDF3) Evaluate(p V3) float64 {
	return sdf_box3d(p, s.Size)
}

func (s *BoxSDF3) BoundingBox() Box3 {
	return Box3{s.Size.Negate(), s.Size}
}

//-----------------------------------------------------------------------------
// 3D Rounded Box

type RoundedBoxSDF3 struct {
	Size   V3
	Radius float64
}

func NewRoundedBoxSDF3(size V3, radius float64) SDF3 {
	// note: store a modified size
	return &RoundedBoxSDF3{size.MulScalar(0.5).SubScalar(radius), radius}
}

func (s *RoundedBoxSDF3) Evaluate(p V3) float64 {
	return sdf_box3d(p, s.Size) - s.Radius
}

func (s *RoundedBoxSDF3) BoundingBox() Box3 {
	d := s.Size.AddScalar(s.Radius)
	return Box3{d.Negate(), d}
}

//-----------------------------------------------------------------------------
// 3D Sphere

type SphereSDF3 struct {
	Radius float64
}

func NewSphereSDF3(radius float64) SDF3 {
	return &SphereSDF3{radius}
}

func (s *SphereSDF3) Evaluate(p V3) float64 {
	return p.Length() - s.Radius
}

func (s *SphereSDF3) BoundingBox() Box3 {
	d := V3{s.Radius, s.Radius, s.Radius}
	return Box3{d.Negate(), d}
}

//-----------------------------------------------------------------------------
// 2D Circle

type CircleSDF2 struct {
	Radius float64
}

func NewCircleSDF2(radius float64) SDF2 {
	return &CircleSDF2{radius}
}

func (s *CircleSDF2) Evaluate(p V2) float64 {
	return p.Length() - s.Radius
}

func (s *CircleSDF2) BoundingBox() Box2 {
	d := V2{s.Radius, s.Radius}
	return Box2{d.Negate(), d}
}

//-----------------------------------------------------------------------------
// 2D Normal Box

type BoxSDF2 struct {
	Size V2
}

func NewBoxSDF2(size V2) SDF2 {
	// note: store a modified size
	return &BoxSDF2{size.MulScalar(0.5)}
}

func (s *BoxSDF2) Evaluate(p V2) float64 {
	return sdf_box2d(p, s.Size)
}

func (s *BoxSDF2) BoundingBox() Box2 {
	return Box2{s.Size.Negate(), s.Size}
}

//-----------------------------------------------------------------------------
// 2D Rounded Box

type RoundedBoxSDF2 struct {
	Size   V2
	Radius float64
}

func NewRoundedBoxSDF2(size V2, radius float64) SDF2 {
	// note: store a modified size
	return &RoundedBoxSDF2{size.MulScalar(0.5).SubScalar(radius), radius}
}

func (s *RoundedBoxSDF2) Evaluate(p V2) float64 {
	return sdf_box2d(p, s.Size) - s.Radius
}

func (s *RoundedBoxSDF2) BoundingBox() Box2 {
	d := s.Size.AddScalar(s.Radius)
	return Box2{d.Negate(), d}
}

//-----------------------------------------------------------------------------
// Transform SDF2

type TransformSDF2 struct {
	Sdf     SDF2
	Matrix  M33
	Inverse M33
}

func NewTransformSDF2(sdf SDF2, matrix M33) SDF2 {
	return &TransformSDF2{sdf, matrix, matrix.Inverse()}
}

func (s *TransformSDF2) Evaluate(p V2) float64 {
	q := s.Inverse.MulPosition(p)
	return s.Sdf.Evaluate(q)
}

func (s *TransformSDF2) BoundingBox() Box2 {
	return s.Matrix.MulBox(s.Sdf.BoundingBox())
}

//-----------------------------------------------------------------------------
// Transform SDF3

type TransformSDF3 struct {
	Sdf     SDF3
	Matrix  M44
	Inverse M44
}

func NewTransformSDF3(sdf SDF3, matrix M44) SDF3 {
	return &TransformSDF3{sdf, matrix, matrix.Inverse()}
}

func (s *TransformSDF3) Evaluate(p V3) float64 {
	q := s.Inverse.MulPosition(p)
	return s.Sdf.Evaluate(q)
}

func (s *TransformSDF3) BoundingBox() Box3 {
	return s.Matrix.MulBox(s.Sdf.BoundingBox())
}

//-----------------------------------------------------------------------------
// Union of SDF3

type UnionSDF3 struct {
	s0  SDF3
	s1  SDF3
	min MinFunc
	k   float64
}

func NewUnionSDF3(s0, s1 SDF3) SDF3 {
	return &UnionSDF3{s0, s1, NormalMin, 0}
}

func NewUnionRoundSDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, RoundMin, k}
}

func NewUnionExpSDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, ExpMin, k}
}

func NewUnionPowSDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, PowMin, k}
}

func NewUnionPolySDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, PolyMin, k}
}

func NewUnionChamferSDF3(s0, s1 SDF3, k float64) SDF3 {
	return &UnionSDF3{s0, s1, ChamferMin, k}
}

func (s *UnionSDF3) Evaluate(p V3) float64 {
	a := s.s0.Evaluate(p)
	b := s.s1.Evaluate(p)
	return s.min(a, b, s.k)
}

func (s *UnionSDF3) BoundingBox() Box3 {
	bb := s.s0.BoundingBox()
	return bb.Extend(s.s1.BoundingBox())
}

//-----------------------------------------------------------------------------