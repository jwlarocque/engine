package r2extra

//

import (
	"math"

	"github.com/golang/geo/r2"
	"github.com/hajimehoshi/ebiten"
)

// == Extra Matrix Transforms ========

// RotatedQuarter returns matrix rotated 90 degrees counter-clockwise
func RotatedQuarter(matrix ebiten.GeoM) ebiten.GeoM {
	ret := ebiten.GeoM{}
	ret.SetElement(0, 0, -matrix.Element(1, 0))
	ret.SetElement(0, 1, -matrix.Element(1, 1))
	ret.SetElement(0, 2, -matrix.Element(1, 2))
	ret.SetElement(1, 0, matrix.Element(0, 0))
	ret.SetElement(1, 1, matrix.Element(0, 1))
	ret.SetElement(1, 2, matrix.Element(0, 2))
	return ret
}

// == Extra Point Functions ========

func ApproxEqual(p, op r2.Point) bool {
	const epsilon = 1e-16
	return math.Abs(p.X-op.X) < epsilon && math.Abs(p.Y-op.Y) < epsilon
}

// ProjectOntoMagnitude returns the magnitude of the projection of p onto op
func ProjectOntoMagnitude(p, op r2.Point) float64 {
	return p.Dot(op) / op.Dot(op)
}

// ProjectOnto returns the projection of p onto op
func ProjectOnto(p, op r2.Point) r2.Point {
	return op.Mul(ProjectOntoMagnitude(p, op))
}
