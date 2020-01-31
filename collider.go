package engine

import (
	"fmt"
)

type Collider struct {
	Vertices             []*Vector2
	center               Vector2 // middle of bounding box
	boundSmall, boundBig Vector2 // bounding box
	Entity
}

// isConvex returns whether the given vertices form a convex polygon
// do not repeat first vertex
// does not handle wack polygons (e.g. self-intersecting)
// FIXME: handle case where two adjacent segments are parallel (i.e., Cross ~= 0) (should always pass)
func isConvex(vertices []*Vector2) bool {
	// triangles, "lines," and "points" are convex
	if len(vertices) <= 3 {
		return true
	}

	// baseline is sign of second to last vector cross last vector
	v1 := vertices[len(vertices)-1].Sub(*vertices[len(vertices)-2])
	v2 := vertices[0].Sub(*vertices[len(vertices)-1])
	positive := v1.Cross(v2) > 0

	// check that sign of all crossZ of adjacent vectors matches baseline
	for i := 0; i < len(vertices)-1; i++ {
		v1 = v2
		v2 = vertices[i+1].Sub(*vertices[i])
		if v1.Cross(v2) > 0 != positive {
			return false
		}
	}
	return true
}

type ErrNotConvex struct {
	ErrStr   string
	Vertices []*Vector2
}

func (e *ErrNotConvex) Error() string {
	return fmt.Sprintf("%s : %v", e.ErrStr, e.Vertices)
}

func NewCollider(vertices []*Vector2) (*Collider, error) {
	if !isConvex(vertices) {
		return nil, &ErrNotConvex{"Collider vertices were not convex.", vertices}
	}
	coll := Collider{}
	coll.Vertices = vertices

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

	coll.center = coll.boundSmall.Add((coll.boundBig.Sub(coll.boundSmall)).Scale(0.5))

	return &coll, nil
}

// GetVertexPos returns the position of the vertex at i plus the collider's position
func (c *Collider) GetVertexPos(i int) Vector2 {
	return c.Vertices[i%len(c.Vertices)].Add(c.Position)
}

func (c *Collider) String() string {
	return fmt.Sprintf("Collider with center: %v, Bounds: (%v, %v), Vertices: %v", c.center, c.boundSmall, c.boundBig, c.Vertices)
}

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

func (c *Collider) satCollides(other *Collider) bool {
	// check each side of this Collider (c)
	var axis Vector2
	var current, cMin, cMax, otherMin, otherMax float64
	for i := 0; i < len(c.Vertices); i++ {
		axis = c.GetVertexPos(i).Sub(c.GetVertexPos(i + 1)).Orthogonal()
		// find projection/"shadow" of c onto axis
		cMin = c.GetVertexPos(0).ProjectOntoMagnitude(axis)
		cMax = cMin
		for j := 1; j < len(c.Vertices); j++ {
			current = c.GetVertexPos(j).ProjectOntoMagnitude(axis)
			if current > cMax {
				cMax = current
			} else if current < cMin {
				cMin = current
			}
		}

		// do the same for other
		otherMin = other.GetVertexPos(0).ProjectOntoMagnitude(axis)
		otherMax = otherMin
		for j := 1; j < len(other.Vertices); j++ {
			current = other.GetVertexPos(j).ProjectOntoMagnitude(axis)
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
