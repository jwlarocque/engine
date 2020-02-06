package r2

import (
	"fmt"
	"math"
)

// == Vector ================
// 2-dimensional Vector stuff
// add functions as necessary

type Vector struct {
	X, Y float64
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector) Angle() float64 {
	return math.Atan2(v.X, v.Y)
}

func (v Vector) String() string {
	return fmt.Sprintf("(%0.24f, %0.24f)", v.X, v.Y)
}

func (v Vector) ApproxEqual(other Vector) bool {
	const epsilon = 1e-16
	return math.Abs(v.X-other.X) < epsilon && math.Abs(v.Y-other.Y) < epsilon
}

func (v Vector) Normalize() Vector {
	len := v.Magnitude()
	if len == 0.0 {
		return Vector{0.0, 0.0}
	}
	return v.Scale(1.0 / len)
}

func (v Vector) Add(other Vector) Vector {
	return Vector{v.X + other.X, v.Y + other.Y}
}

func (v Vector) Sub(other Vector) Vector {
	return Vector{v.X - other.X, v.Y - other.Y}
}

func (v Vector) Scale(scalar float64) Vector {
	return Vector{v.X * scalar, v.Y * scalar}
}

func (v Vector) Dot(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vector) Cross(other Vector) float64 {
	return v.X*other.Y - other.X*v.Y
}

func (v Vector) Orthogonal() Vector {
	return Vector{v.Y, -v.X}
}

// ProjectOntoMagnitude returns the magnitude of the projection of v onto other
func (v Vector) ProjectOntoMagnitude(other Vector) float64 {
	return v.Dot(other) / other.Dot(other)
}

// ProjectOnto returns the projection of v onto other
func (v Vector) ProjectOnto(other Vector) Vector {
	return other.Scale(v.ProjectOntoMagnitude(other))
}
