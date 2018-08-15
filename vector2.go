package main

import (
	"math"
)

// Vector2 is the struct that represents a vector in 2D space
type Vector2 struct {
	X float64
	Y float64
}

// Add is the function that performances two vectors addition (v + u)
func (v Vector2) Add(u Vector2) Vector2 {
	return Vector2{
		X: v.X + u.X,
		Y: v.Y + u.Y,
	}
}

// Sub is the function that performances two vectors subtraction (v - u)
func (v Vector2) Sub(u Vector2) Vector2 {
	return Vector2{
		X: v.X - u.X,
		Y: v.Y - u.Y,
	}
}

// Multiply is the function that performances vector scale
func (v Vector2) Multiply(u float64) Vector2 {
	return Vector2{
		X: v.X * u,
		Y: v.Y * u,
	}
}

// Magnitude is the function that return the magnitude of this vector
func (v Vector2) Magnitude() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

// Normalize is the function that return the norm form of this vector
func (v Vector2) Normalize() Vector2 {
	length := v.Magnitude()
	return Vector2{
		v.X / length,
		v.Y / length,
	}
}

// Dot is the function that return the dot product of v and u
func (v Vector2) Dot(u Vector2) float64 {
	return v.X*u.X + v.Y*u.Y
}
