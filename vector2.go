package engine

// 2-dimensional vector stuff
// add functions as necessary

import (
	"fmt"
	"math"
)

type Vector2 struct {
	X, Y float64
}

func (v Vector2) Magnitude() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

func (v Vector2) Angle() float64 {
	return math.Atan2(v.X, v.Y)
}

func (v Vector2) String() string {
	return fmt.Sprintf("(%0.24f, %0.24f)", v.X, v.Y)
}

func (v Vector2) ApproxEqual(other Vector2) bool {
	const epsilon = 1e-16
	return math.Abs(v.X-other.X) < epsilon && math.Abs(v.Y-other.Y) < epsilon
}

func (v Vector2) Normalize() Vector2 {
	len := v.Magnitude()
	if len == 0 {
		return Vector2{0, 0}
	}
	return v.Scale(1 / math.Sqrt(len))
}

func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{v.X + other.X, v.Y + other.Y}
}

func (v Vector2) Scale(scalar float64) Vector2 {
	return Vector2{v.X * scalar, v.Y * scalar}
}
