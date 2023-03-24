//-----------------------------------------------------------------------------
/*

2D Signed Distance Functions

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"errors"
	"math"

	"github.com/gmlewis/sdfx/vec/conv"
	"github.com/gmlewis/sdfx/vec/p2"
	v2 "github.com/gmlewis/sdfx/vec/v2"
	"github.com/gmlewis/sdfx/vec/v2i"
	v3 "github.com/gmlewis/sdfx/vec/v3"
)

//-----------------------------------------------------------------------------

// SDF2 is the interface to a 2d signed distance function object.
type SDF2 interface {
	Evaluate(p v2.Vec) float64
	BoundingBox() Box2
}

//-----------------------------------------------------------------------------
// SDF2 Evaluation Caching (experimental)

type sdf2Cache struct {
	cache map[v2.Vec]float64
	hits  uint
}

func (c *sdf2Cache) lookup(p v2.Vec) (float64, error) {
	if d, ok := c.cache[p]; ok {
		c.hits++
		return d, nil
	}
	return 0, errors.New("not found")
}

func (c *sdf2Cache) store(p v2.Vec, d float64) {
	c.cache[p] = d
}

func newSdf2Cache() *sdf2Cache {
	c := sdf2Cache{}
	c.cache = make(map[v2.Vec]float64)
	return &c
}

//-----------------------------------------------------------------------------
// Basic SDF Functions

func sdfBox2d(p, s v2.Vec) float64 {
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

//-----------------------------------------------------------------------------
// 2D Circle

// CircleSDF2 is the 2d signed distance object for a circle.
type CircleSDF2 struct {
	radius float64
	bb     Box2
}

// Circle2D returns the SDF2 for a 2d circle.
func Circle2D(radius float64) (SDF2, error) {
	if radius < 0 {
		return nil, ErrMsg("radius < 0")
	}
	s := CircleSDF2{}
	s.radius = radius
	d := v2.Vec{radius, radius}
	s.bb = Box2{d.Neg(), d}
	return &s, nil
}

// Evaluate returns the minimum distance to a 2d circle.
func (s *CircleSDF2) Evaluate(p v2.Vec) float64 {
	return p.Length() - s.radius
}

// BoundingBox returns the bounding box of a 2d circle.
func (s *CircleSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
// 2D Box (rounded corners with round > 0)

// BoxSDF2 is the 2d signed distance object for a rectangular box.
type BoxSDF2 struct {
	size  v2.Vec
	round float64
	bb    Box2
}

// Box2D returns a 2d box.
func Box2D(size v2.Vec, round float64) SDF2 {
	size = size.MulScalar(0.5)
	s := BoxSDF2{}
	s.size = size.SubScalar(round)
	s.round = round
	s.bb = Box2{size.Neg(), size}
	return &s
}

// Evaluate returns the minimum distance to a 2d box.
func (s *BoxSDF2) Evaluate(p v2.Vec) float64 {
	return sdfBox2d(p, s.size) - s.round
}

// BoundingBox returns the bounding box for a 2d box.
func (s *BoxSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
// 2D Line

// LineSDF2 is the 2d signed distance object for a line.
type LineSDF2 struct {
	l     float64 // line length
	round float64 // rounding
	bb    Box2    // bounding box
}

// Line2D returns a line from (-l/2,0) to (l/2,0).
func Line2D(l, round float64) SDF2 {
	s := LineSDF2{}
	s.l = l / 2
	s.round = round
	s.bb = Box2{v2.Vec{-s.l - round, -round}, v2.Vec{s.l + round, round}}
	return &s
}

// Evaluate returns the minimum distance to a 2d line.
func (s *LineSDF2) Evaluate(p v2.Vec) float64 {
	p = p.Abs()
	if p.X <= s.l {
		return p.Y - s.round
	}
	return p.Sub(v2.Vec{s.l, 0}).Length() - s.round
}

// BoundingBox returns the bounding box for a 2d line.
func (s *LineSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
// 2D Spiral

type V2 = v2.Vec

type SpiralSDF2 struct {
	start float64 // start angle (and radius) in radians
	end   float64 // end angle (and radius) in radians
	sθ    float64 // start normalized angle (-pi <= sθ <= pi)
	eθ    float64 // end normalized angle (pi <= eθ <= pi)
	ps    V2      // start point in cartesian coordinates
	pe    V2      // end point in cartesian coordinates
	round float64 // rounding
	bb    Box2    // bounding box
}

// Spiral2D returns a spiral with the equation `r = θ`
// starting at radius (and angle) `start` in radians
// and ending at radius (and angle) `end` in radians.
// `start` must be less than or equal to `end`.
func Spiral2D(start, end, round float64) SDF2 {
	ps := V2{X: start * math.Cos(start), Y: start * math.Sin(start)}
	pe := V2{X: end * math.Cos(end), Y: end * math.Sin(end)}
	return &SpiralSDF2{
		start: start,
		end:   end,
		sθ:    math.Atan2(ps.Y, ps.X),
		eθ:    math.Atan2(pe.Y, pe.X),
		ps:    ps,
		pe:    pe,
		round: round,
		bb:    Box2{V2{-end - round, -end - round}, V2{end + round, end + round}},
	}
}

// Evaluate returns the minimum distance to the spiral.
func (s *SpiralSDF2) Evaluate(p V2) float64 {
	pr := p.Length()
	pθ := math.Atan2(p.Y, p.X)
	c := 1 - math.Cos(pθ-pr)
	dist := 0.5 * math.Pi * math.Sqrt(c)

	ds := s.ps.Sub(p).Length()
	if ds < dist {
		dist = ds
	}
	de := s.pe.Sub(p).Length()
	if de < dist {
		dist = de
	}

	if s.start > 0 && pr < s.start+math.Pi {
		dist = ds
		if de < dist {
			dist = de
		}
		delta := pθ - s.sθ
		for delta < 0 {
			delta += 2 * math.Pi
		}
		angle := s.start + delta
		if angle <= s.end {
			sp := V2{X: angle * math.Cos(angle), Y: angle * math.Sin(angle)}
			d := sp.Sub(p).Length()
			if d < dist {
				dist = d
			}
		}
	} else if pr > s.end-math.Pi {
		dist = de
		if ds < dist {
			dist = ds
		}
		delta := pθ - s.eθ
		for delta > 0 {
			delta -= 2 * math.Pi
		}
		angle := s.end + delta
		if angle >= s.start {
			sp := V2{X: angle * math.Cos(angle), Y: angle * math.Sin(angle)}
			d := sp.Sub(p).Length()
			if d < dist {
				dist = d
			}
		}
	}
	return dist - s.round
}

// BoundingBox returns the bounding box for the spiral.
func (s *SpiralSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// OffsetSDF2 offsets the distance function of an existing SDF2.
type OffsetSDF2 struct {
	sdf    SDF2
	offset float64
	bb     Box2
}

// Offset2D returns an SDF2 that offsets the distance function of another SDF2.
func Offset2D(sdf SDF2, offset float64) SDF2 {
	s := OffsetSDF2{}
	s.sdf = sdf
	s.offset = offset
	// work out the bounding box
	bb := sdf.BoundingBox()
	s.bb = NewBox2(bb.Center(), bb.Size().AddScalar(2*offset))
	return &s
}

// Evaluate returns the minimum distance to an offset SDF2.
func (s *OffsetSDF2) Evaluate(p v2.Vec) float64 {
	return s.sdf.Evaluate(p) - s.offset
}

// BoundingBox returns the bounding box of an offset SDF2.
func (s *OffsetSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// IntersectionSDF2 is the intersection of two SDF2s.
type IntersectionSDF2 struct {
	s0  SDF2
	s1  SDF2
	max MaxFunc
	bb  Box2
}

// Intersect2D returns the intersection of two SDF2s.
func Intersect2D(s0, s1 SDF2) SDF2 {
	if s0 == nil || s1 == nil {
		return nil
	}
	s := IntersectionSDF2{}
	s.s0 = s0
	s.s1 = s1
	s.max = math.Max
	// TODO fix bounding box
	s.bb = s0.BoundingBox()
	return &s
}

// Evaluate returns the minimum distance to the SDF2 intersection.
func (s *IntersectionSDF2) Evaluate(p v2.Vec) float64 {
	return s.max(s.s0.Evaluate(p), s.s1.Evaluate(p))
}

// SetMax sets the maximum function to control blending.
func (s *IntersectionSDF2) SetMax(max MaxFunc) {
	s.max = max
}

// BoundingBox returns the bounding box of an SDF2 intersection.
func (s *IntersectionSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Cut an SDF2 along a line

// CutSDF2 is an SDF2 made by cutting across an existing SDF2.
type CutSDF2 struct {
	sdf SDF2
	a   v2.Vec // point on line
	n   v2.Vec // normal to line
	bb  Box2   // bounding box
}

// Cut2D cuts the SDF2 along a line from a in direction v.
// The SDF2 to the right of the line remains.
func Cut2D(sdf SDF2, a, v v2.Vec) SDF2 {
	s := CutSDF2{}
	s.sdf = sdf
	s.a = a
	v = v.Normalize()
	s.n = v2.Vec{-v.Y, v.X}
	// TODO - cut the bounding box
	s.bb = sdf.BoundingBox()
	return &s
}

// Evaluate returns the minimum distance to cut SDF2.
func (s *CutSDF2) Evaluate(p v2.Vec) float64 {
	return math.Max(p.Sub(s.a).Dot(s.n), s.sdf.Evaluate(p))
}

// BoundingBox returns the bounding box for the cut SDF2.
func (s *CutSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Transform SDF2 (rotation and translation are distance preserving)

// TransformSDF2 transorms an SDF2 with rotation, translation and scaling.
type TransformSDF2 struct {
	sdf  SDF2
	mInv M33
	bb   Box2
}

// Transform2D applies a transformation matrix to an SDF2.
// Distance is *not* preserved with scaling.
func Transform2D(sdf SDF2, m M33) SDF2 {
	s := TransformSDF2{}
	s.sdf = sdf
	s.mInv = m.Inverse()
	s.bb = m.MulBox(sdf.BoundingBox())
	return &s
}

// Evaluate returns the minimum distance to a transformed SDF2.
// Distance is *not* preserved with scaling.
func (s *TransformSDF2) Evaluate(p v2.Vec) float64 {
	q := s.mInv.MulPosition(p)
	return s.sdf.Evaluate(q)
}

// BoundingBox returns the bounding box of a transformed SDF2.
func (s *TransformSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------
// Uniform XY Scaling of SDF2s (we can work out the distance)

// ScaleUniformSDF2 scales another SDF2 on each axis.
type ScaleUniformSDF2 struct {
	sdf     SDF2
	k, invk float64
	bb      Box2
}

// ScaleUniform2D scales an SDF2 by k on each axis.
// Distance is correct with scaling.
func ScaleUniform2D(sdf SDF2, k float64) SDF2 {
	m := Scale2d(v2.Vec{k, k})
	return &ScaleUniformSDF2{
		sdf:  sdf,
		k:    k,
		invk: 1.0 / k,
		bb:   m.MulBox(sdf.BoundingBox()),
	}
}

// Evaluate returns the minimum distance to an SDF2 with uniform scaling.
func (s *ScaleUniformSDF2) Evaluate(p v2.Vec) float64 {
	q := p.MulScalar(s.invk)
	return s.sdf.Evaluate(q) * s.k
}

// BoundingBox returns the bounding box of an SDF2 with uniform scaling.
func (s *ScaleUniformSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// Center2D centers the origin of an SDF2 on it's bounding box.
func Center2D(s SDF2) SDF2 {
	ofs := s.BoundingBox().Center().Neg()
	return Transform2D(s, Translate2d(ofs))
}

// CenterAndScale2D centers the origin of an SDF2 on it's bounding box, and then scales it.
// Distance is correct with scaling.
func CenterAndScale2D(s SDF2, k float64) SDF2 {
	ofs := s.BoundingBox().Center().Neg()
	s = Transform2D(s, Translate2d(ofs))
	return ScaleUniform2D(s, k)
}

//-----------------------------------------------------------------------------
// ArraySDF2: Create an X by Y array of a given SDF2

// ArraySDF2 defines an XY grid array of an existing SDF2.
type ArraySDF2 struct {
	sdf  SDF2
	num  v2i.Vec // grid size
	step v2.Vec  // grid step size
	min  MinFunc
	bb   Box2
}

// Array2D returns an XY grid array of an existing SDF2.
func Array2D(sdf SDF2, num v2i.Vec, step v2.Vec) SDF2 {
	// check the number of steps
	if num.X <= 0 || num.Y <= 0 {
		return nil
	}
	s := ArraySDF2{}
	s.sdf = sdf
	s.num = num
	s.step = step
	s.min = math.Min
	// work out the bounding box
	bb0 := sdf.BoundingBox()
	bb1 := bb0.Translate(step.Mul(conv.V2iToV2(num.SubScalar(1))))
	s.bb = bb0.Extend(bb1)
	return &s
}

// SetMin sets the minimum function to control blending.
func (s *ArraySDF2) SetMin(min MinFunc) {
	s.min = min
}

// Evaluate returns the minimum distance to a grid array of SDF2s.
func (s *ArraySDF2) Evaluate(p v2.Vec) float64 {
	d := math.MaxFloat64
	for j := 0; j < s.num.X; j++ {
		for k := 0; k < s.num.Y; k++ {
			x := p.Sub(v2.Vec{float64(j) * s.step.X, float64(k) * s.step.Y})
			d = s.min(d, s.sdf.Evaluate(x))
		}
	}
	return d
}

// BoundingBox returns the bounding box of a grid array of SDF2s.
func (s *ArraySDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// RotateUnionSDF2 defines a union of rotated SDF2s.
type RotateUnionSDF2 struct {
	sdf  SDF2
	num  int
	step M33
	min  MinFunc
	bb   Box2
}

// RotateUnion2D returns a union of rotated SDF2s.
func RotateUnion2D(sdf SDF2, num int, step M33) SDF2 {
	// check the number of steps
	if num <= 0 {
		return nil
	}
	s := RotateUnionSDF2{}
	s.sdf = sdf
	s.num = num
	s.step = step.Inverse()
	s.min = math.Min
	// work out the bounding box
	v := sdf.BoundingBox().Vertices()
	bbMin := v[0]
	bbMax := v[0]
	for i := 0; i < s.num; i++ {
		bbMin = bbMin.Min(v.Min())
		bbMax = bbMax.Max(v.Max())
		mulVertices2(v, step)
	}
	s.bb = Box2{bbMin, bbMax}
	return &s
}

// Evaluate returns the minimum distance to a union of rotated SDF2s.
func (s *RotateUnionSDF2) Evaluate(p v2.Vec) float64 {
	d := math.MaxFloat64
	rot := Identity2d()
	for i := 0; i < s.num; i++ {
		x := rot.MulPosition(p)
		d = s.min(d, s.sdf.Evaluate(x))
		rot = rot.Mul(s.step)
	}
	return d
}

// SetMin sets the minimum function to control blending.
func (s *RotateUnionSDF2) SetMin(min MinFunc) {
	s.min = min
}

// BoundingBox returns the bounding box of a union of rotated SDF2s.
func (s *RotateUnionSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// RotateCopySDF2 copies an SDF2 n times in a full circle.
type RotateCopySDF2 struct {
	sdf   SDF2
	theta float64
	bb    Box2
}

// RotateCopy2D rotates and copies an SDF2 n times in a full circle.
func RotateCopy2D(sdf SDF2, n int) SDF2 {
	// check the number of steps
	if n <= 0 {
		return nil
	}
	s := RotateCopySDF2{}
	s.sdf = sdf
	s.theta = Tau / float64(n)
	// work out the bounding box
	bb := sdf.BoundingBox()
	rmax := 0.0
	// find the bounding box vertex with the greatest distance from the origin
	for _, v := range bb.Vertices() {
		l := v.Length()
		if l > rmax {
			rmax = l
		}
	}
	s.bb = Box2{v2.Vec{-rmax, -rmax}, v2.Vec{rmax, rmax}}
	return &s
}

// Evaluate returns the minimum distance to a rotate/copy SDF2.
func (s *RotateCopySDF2) Evaluate(p v2.Vec) float64 {
	// Map p to a point in the first copy sector.
	pnew := conv.P2ToV2(p2.Vec{p.Length(), SawTooth(math.Atan2(p.Y, p.X), s.theta)})
	return s.sdf.Evaluate(pnew)
}

// BoundingBox returns the bounding box of a rotate/copy SDF2.
func (s *RotateCopySDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// SliceSDF2 creates an SDF2 from a planar slice through an SDF3.
type SliceSDF2 struct {
	sdf SDF3   // the sdf3 being sliced
	a   v3.Vec // 3d point for 2d origin
	u   v3.Vec // vector for the 2d x-axis
	v   v3.Vec // vector for the 2d y-axis
	bb  Box2   // bounding box
}

// Slice2D returns an SDF2 created from a planar slice through an SDF3.
func Slice2D(
	sdf SDF3, // SDF3 to be sliced
	a v3.Vec, // point on slicing plane
	n v3.Vec, // normal to slicing plane
) SDF2 {
	s := SliceSDF2{}
	s.sdf = sdf
	s.a = a
	// work out the x/y vectors on the plane.
	if n.X == 0 {
		s.u = v3.Vec{1, 0, 0}
	} else if n.Y == 0 {
		s.u = v3.Vec{0, 1, 0}
	} else if n.Z == 0 {
		s.u = v3.Vec{0, 0, 1}
	} else {
		s.u = v3.Vec{n.Y, -n.X, 0}
	}
	s.v = n.Cross(s.u)
	s.u = s.u.Normalize()
	s.v = s.v.Normalize()
	// work out the bounding box
	// TODO: This is bigger than it needs to be. We could consider intersection
	// between the plane and the edges of the 3d bounding box for a smaller 2d
	// bounding box in some circumstances.
	v3Verts := sdf.BoundingBox().Vertices()
	v2Verts := make(v2.VecSet, len(v3Verts))
	n = n.Normalize()
	for i, v := range v3Verts {
		// project the 3d bounding box vertex onto the plane
		va := v.Sub(s.a)
		pa := va.Sub(n.MulScalar(n.Dot(va)))
		// work out the 3d point in terms of the 2d unit vectors
		v2Verts[i] = v2.Vec{pa.Dot(s.u), pa.Dot(s.v)}
	}
	s.bb = Box2{v2Verts.Min(), v2Verts.Max()}
	return &s
}

// Evaluate returns the minimum distance to the sliced SDF2.
func (s *SliceSDF2) Evaluate(p v2.Vec) float64 {
	pnew := s.a.Add(s.u.MulScalar(p.X)).Add(s.v.MulScalar(p.Y))
	return s.sdf.Evaluate(pnew)
}

// BoundingBox returns the bounding box of the sliced SDF2.
func (s *SliceSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// UnionSDF2 is a union of multiple SDF2 objects.
type UnionSDF2 struct {
	sdf []SDF2
	min MinFunc
	bb  Box2
}

// Union2D returns the union of multiple SDF2 objects.
func Union2D(sdf ...SDF2) SDF2 {
	if len(sdf) == 0 {
		return nil
	}
	s := UnionSDF2{}
	// strip out any nils
	s.sdf = make([]SDF2, 0, len(sdf))
	for _, x := range sdf {
		if x != nil {
			s.sdf = append(s.sdf, x)
		}
	}
	if len(s.sdf) == 0 {
		return nil
	}
	if len(s.sdf) == 1 {
		// only one sdf - not really a union
		return s.sdf[0]
	}
	// work out the bounding box
	bb := s.sdf[0].BoundingBox()
	for _, x := range s.sdf {
		bb = bb.Extend(x.BoundingBox())
	}
	s.bb = bb
	s.min = math.Min
	return &s
}

// Evaluate returns the minimum distance to the SDF2 union.
func (s *UnionSDF2) Evaluate(p v2.Vec) float64 {

	// work out the min/max distance for every bounding box
	vs := make([]v2.Vec, len(s.sdf))
	minDist2 := -1.0
	minIndex := 0
	for i := range s.sdf {
		vs[i] = s.sdf[i].BoundingBox().MinMaxDist2(p)
		// as we go record the sdf with the minimum minimum d2 value
		if minDist2 < 0 || vs[i].X < minDist2 {
			minDist2 = vs[i].X
			minIndex = i
		}
	}

	var d float64
	first := true
	for i := range s.sdf {
		// only an sdf whose min/max distances overlap
		// the minimum box are worthy of consideration
		if i == minIndex || vs[minIndex].Overlap(vs[i]) {
			x := s.sdf[i].Evaluate(p)
			if first {
				first = false
				d = x
			} else {
				d = s.min(d, x)
			}
		}
	}
	return d
}

// EvaluateSlow returns the minimum distance to the SDF2 union.
func (s *UnionSDF2) EvaluateSlow(p v2.Vec) float64 {
	var d float64
	for i := range s.sdf {
		x := s.sdf[i].Evaluate(p)
		if i == 0 {
			d = x
		} else {
			d = s.min(d, x)
		}
	}
	return d
}

// SetMin sets the minimum function to control SDF2 blending.
func (s *UnionSDF2) SetMin(min MinFunc) {
	s.min = min
}

// BoundingBox returns the bounding box of an SDF2 union.
func (s *UnionSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// DifferenceSDF2 is the difference of two SDF2s.
type DifferenceSDF2 struct {
	s0  SDF2
	s1  SDF2
	max MaxFunc
	bb  Box2
}

// Difference2D returns the difference of two SDF2 objects, s0 - s1.
func Difference2D(s0, s1 SDF2) SDF2 {
	if s1 == nil {
		return s0
	}
	if s0 == nil {
		return nil
	}
	s := DifferenceSDF2{}
	s.s0 = s0
	s.s1 = s1
	s.max = math.Max
	s.bb = s0.BoundingBox()
	return &s
}

// Evaluate returns the minimum distance to the difference of two SDF2s.
func (s *DifferenceSDF2) Evaluate(p v2.Vec) float64 {
	return s.max(s.s0.Evaluate(p), -s.s1.Evaluate(p))
}

// SetMax sets the maximum function to control blending.
func (s *DifferenceSDF2) SetMax(max MaxFunc) {
	s.max = max
}

// BoundingBox returns the bounding box of the difference of two SDF2s.
func (s *DifferenceSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// ElongateSDF2 is the elongation of an SDF2.
type ElongateSDF2 struct {
	sdf    SDF2   // the sdf being elongated
	hp, hn v2.Vec // positive/negative elongation vector
	bb     Box2   // bounding box
}

// Elongate2D returns the elongation of an SDF2.
func Elongate2D(sdf SDF2, h v2.Vec) SDF2 {
	h = h.Abs()
	s := ElongateSDF2{
		sdf: sdf,
		hp:  h.MulScalar(0.5),
		hn:  h.MulScalar(-0.5),
	}
	// bounding box
	bb := sdf.BoundingBox()
	bb0 := bb.Translate(s.hp)
	bb1 := bb.Translate(s.hn)
	s.bb = bb0.Extend(bb1)
	return &s
}

// Evaluate returns the minimum distance to an elongated SDF2.
func (s *ElongateSDF2) Evaluate(p v2.Vec) float64 {
	q := p.Sub(p.Clamp(s.hn, s.hp))
	return s.sdf.Evaluate(q)
}

// BoundingBox returns the bounding box of an elongated SDF2.
func (s *ElongateSDF2) BoundingBox() Box2 {
	return s.bb
}

//-----------------------------------------------------------------------------

// GenerateMesh2D generates a set of internal mesh points for an SDF2.
func GenerateMesh2D(s SDF2, grid v2i.Vec) (v2.VecSet, error) {

	// create the grid mapping for the bounding box
	m, err := NewMap2(s.BoundingBox(), grid, false)
	if err != nil {
		return nil, err
	}

	// create the vertex set storage
	vset := make(v2.VecSet, 0, grid.X*grid.Y)

	// iterate across the grid and add the vertices if they are inside the SDF2
	for i := 0; i < grid.X; i++ {
		for j := 0; j < grid.Y; j++ {
			v := m.ToV2(v2i.Vec{i, j})
			if s.Evaluate(v) <= 0 {
				vset = append(vset, v)
			}
		}
	}

	return vset, nil
}

//-----------------------------------------------------------------------------

// LineOf2D returns a union of 2D objects positioned along a line from p0 to p1.
func LineOf2D(s SDF2, p0, p1 v2.Vec, pattern string) SDF2 {
	var objects []SDF2
	if pattern != "" {
		x := p0
		dx := p1.Sub(p0).DivScalar(float64(len(pattern)))
		for _, c := range pattern {
			if c == 'x' {
				objects = append(objects, Transform2D(s, Translate2d(x)))
			}
			x = x.Add(dx)
		}
	}
	return Union2D(objects...)
}

//-----------------------------------------------------------------------------

// Multi2D creates a union of an SDF2 at a set of 2D positions.
func Multi2D(s SDF2, positions v2.VecSet) SDF2 {
	if (s == nil) || (len(positions) == 0) {
		return nil
	}
	objects := make([]SDF2, len(positions))
	for i, p := range positions {
		objects[i] = Transform2D(s, Translate2d(p))
	}
	return Union2D(objects...)
}

//-----------------------------------------------------------------------------
