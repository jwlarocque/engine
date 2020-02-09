/* TODO: Convert to interface
 *       Add point, circle, rect, etc. colliders with optimized detection
 *       Consider removing center (of bounding box) - it's currently unused.
 */

package mech

import (
	"fmt"

	"github.com/golang/geo/r2"
	"github.com/jwlarocque/engine/r2extra"
)

// Collider has Vertices and an Entity to keep track of its position in the level.
// It can determine whether it is intersecting/overlapping with another Collider.
// Note: Vertices must form a convex polygon (do not repeat first/last vertex).
type Collider struct {
	Vertices             []*r2.Point
	center               r2.Point // middle of bounding box
	boundSmall, boundBig r2.Point // bounding box
	Position             r2.Point
	Velocity             r2.Point
}

// isConvex returns whether the given vertices form a convex polygon
// does not handle wack polygons (e.g. self-intersecting)
// FIXME: handle case where two adjacent segments are parallel (i.e., Cross ~= 0) (should always pass)
// FIXME: I don't know how Collides will handle colliders with fewer than 4 vertices
func isConvex(vertices []*r2.Point) bool {
	// triangles, "lines," and "points" are convex
	if len(vertices) <= 3 {
		return true
	}

	// baseline is sign of second to last vector cross last vector
	v1 := vertices[len(vertices)-1].Sub(*vertices[len(vertices)-2])
	v2 := vertices[0].Sub(*vertices[len(vertices)-1])
	positive := v1.Cross(v2) > 0

	// check that sign of all Cross of adjacent vectors matches baseline
	for i := 0; i < len(vertices)-1; i++ {
		v1 = v2
		v2 = vertices[i+1].Sub(*vertices[i])
		if v1.Cross(v2) > 0 != positive {
			return false
		}
	}
	return true
}

// ErrNotConvex is returned by NewCollider when the provided vertices are not convex
type ErrNotConvex struct {
	ErrStr   string
	Vertices []*r2.Point
}

func (e *ErrNotConvex) Error() string {
	return fmt.Sprintf("%s : %v", e.ErrStr, e.Vertices)
}

// NewCollider constructs a new Collider from the provided vertices
func NewCollider(vertices []*r2.Point) (*Collider, error) {
	if !isConvex(vertices) {
		return nil, &ErrNotConvex{"Collider vertices were not convex.", vertices}
	}
	coll := Collider{}
	coll.Vertices = vertices

	// find bounding box
	coll.boundSmall = *vertices[0]
	coll.boundBig = *vertices[0]
	for i := 1; i < len(vertices); i++ {
		if vertices[i].X < coll.boundSmall.X {
			coll.boundSmall.X = vertices[i].X
		} else if vertices[i].X > coll.boundBig.X {
			coll.boundBig.X = vertices[i].X
		}
		if vertices[i].Y < coll.boundSmall.Y {
			coll.boundSmall.Y = vertices[i].Y
		} else if vertices[i].Y > coll.boundBig.Y {
			coll.boundBig.Y = vertices[i].Y
		}
	}

	// find center of bounding box
	coll.center = coll.boundSmall.Add((coll.boundBig.Sub(coll.boundSmall)).Mul(0.5))
	return &coll, nil
}

// GetVertexPos returns the position of the vertex at i (% len(vertices))
// plus the Collider's position
func (c *Collider) GetVertexPos(i int) r2.Point {
	return c.Vertices[i%len(c.Vertices)].Add(c.Position)
}

func (c *Collider) String() string {
	return fmt.Sprintf("Collider with center: %v, Bounds: (%v, %v), Vertices: %v", c.center, c.boundSmall, c.boundBig, c.Vertices)
}

//
// == Time Blur ========

// Blur "streches" a collider from its current position to its
// position after time
// !!! Creates a new collider.  Use sparingly. !!!
func (c Collider) Blur(time float64) Collider {
	return Collider{}
}

//
// == Collision Detection ========

// checks for bounding box collision
func (c *Collider) bBoxCollides(other *Collider) bool {
	// TODO: this is dumb, write cleaner way
	cS := c.boundSmall.Add(c.Position)
	cB := c.boundBig.Add(c.Position)
	oS := other.boundSmall.Add(other.Position)
	oB := other.boundBig.Add(other.Position)
	if cS.X > oB.X || cS.Y > oB.Y || cB.X < oS.X || cB.Y < oS.Y {
		return false
	}
	return true
}

// check for polygon collision with separating axis theorem
// idea: https://www.sevenson.com.au/actionscript/sat/
// TODO: Only checks edges of c as separating line, so must
//       && with other.satCollides(c) - combine into single call
func (c *Collider) satCollides(other *Collider) bool {
	// check each side of this Collider (c)
	var axis r2.Point
	var current, cMin, cMax, otherMin, otherMax float64
	for i := 0; i < len(c.Vertices); i++ {
		// axis is ortho to a side of c
		axis = c.GetVertexPos(i).Sub(c.GetVertexPos(i + 1)).Ortho()
		// find projection/"shadow" of c onto axis
		cMin = r2extra.ProjectOntoMagnitude(c.GetVertexPos(0), axis)
		cMax = cMin
		for j := 1; j < len(c.Vertices); j++ {
			current = r2extra.ProjectOntoMagnitude(c.GetVertexPos(j), axis)
			if current > cMax {
				cMax = current
			} else if current < cMin {
				cMin = current
			}
		}

		// do the same for other
		otherMin = r2extra.ProjectOntoMagnitude(other.GetVertexPos(0), axis)
		otherMax = otherMin
		for j := 1; j < len(other.Vertices); j++ {
			current = r2extra.ProjectOntoMagnitude(other.GetVertexPos(j), axis)
			if current > otherMax {
				otherMax = current
			} else if current < otherMin {
				otherMin = current
			}
		}

		// if projections don't overlap, there must be a gap between the colliders
		if cMax < otherMin || otherMax < cMin {
			return false
		}
	}
	return true
}

// Collides returns whether this Collider intersects with the other Collider
func (c *Collider) Collides(other *Collider) bool {
	// check bounding box intersection
	if !c.bBoxCollides(other) {
		return false
	}
	// check for gaps with separating axis theorem
	// note: satCollides only considers the axes of c, so just and the results together to consider both
	return c.satCollides(other) && other.satCollides(c)
}

// WillCollide returns whether c and other _are_ colliding after timeSteps
// assumes collider moves exactly velocity every time step (i.e., a time step is one unit of time)
func (c Collider) WillCollide(other Collider, timeSteps int) bool {
	c.Position.Add(c.Velocity)
	other.Position.Add(other.Velocity)
	return c.Collides(&other)
}
