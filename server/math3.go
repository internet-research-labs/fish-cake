package server

import (
	"math"
)

// Vector3 represnts the normal 3d cartesian -vector
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Add returns a new vector3
func Add(a, b Vector3) Vector3 {
	return Vector3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Add returns a new vector3
func Sub(a, b Vector3) Vector3 {
	return Vector3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Scale returns a vector scaled by a constant
func Scale(p Vector3, a float64) Vector3 {
	return Vector3{a * p.X, a * p.Y, a * p.Z}
}

// Return norm of a point
func Norm(p Vector3) float64 {
	return p.X*p.X + p.Y*p.Y + p.Z*p.Z
}

// Return distance between 2 points
func Distance(p, q Vector3) float64 {
	v := Sub(p, q)
	return math.Max(v.X, math.Max(v.Y, v.Z))
}
