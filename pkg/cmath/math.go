package cmath

import (
	"math"
)

type Vec3 struct {
	X float32 `json:"x" db:"x"`
	Y float32 `json:"y" db:"y"`
	Z float32 `json:"z" db:"z"`
}

// TODO: rename to "ObjectPosition"
type ObjectPosition struct {
	Location Vec3 `db:"location" json:"location"`
	Rotation Vec3 `db:"rotation" json:"rotation"`
	Scale    Vec3 `db:"scale" json:"scale"`
}

type UserPosition struct {
	Location Vec3 `db:"location" json:"location"`
	Rotation Vec3 `db:"rotation" json:"rotation"`
}

func (v *Vec3) Plus(v2 Vec3) {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
}

func Add(v1 Vec3, v2 Vec3) Vec3 {
	return Vec3{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func MultiplyN(v Vec3, n float32) Vec3 {
	return Vec3{
		X: v.X * n,
		Y: v.Y * n,
		Z: v.Z * n,
	}
}

func (v *Vec3) ToVec3f64() Vec3f64 {
	return Vec3f64{
		float64(v.X),
		float64(v.Y),
		float64(v.Z),
	}
}

type Vec3f64 struct {
	X float64 `json:"x" db:"x"`
	Y float64 `json:"y" db:"y"`
	Z float64 `json:"z" db:"z"`
}

func (v *Vec3f64) ToVec3() Vec3 {
	return Vec3{
		float32(v.X),
		float32(v.Y),
		float32(v.Z),
	}
}

func MNan32Vec3() Vec3 {
	return Vec3{ // What is this, TODO: make method on math package
		X: float32(math.NaN()),
		Y: float32(math.NaN()),
		Z: float32(math.NaN()),
	}
}

func MNaN32() float32 {
	return float32(math.NaN())
}

func Distance(x, y *Vec3) float64 {
	r := 0.0
	q := float64(x.X - y.X)
	r += q * q
	q = float64(x.Y - y.Y)
	r += q * q
	q = float64(x.Z - y.Z)
	r += q * q
	return math.Sqrt(r)
}

// DefaultPosition FIXME: Magic numbers ?
func DefaultPosition() Vec3 {
	return Vec3{
		X: -170.0,
		Y: 50.0,
		Z: 340.0,
	}
}
