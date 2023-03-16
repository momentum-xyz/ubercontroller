//go:generate go run gen/mus.go

package cmath

import (
	"encoding/binary"
	"math"
)

const (
	Float32Bytes = 4
)

type Vec3 struct {
	X float32 `json:"x" db:"x"`
	Y float32 `json:"y" db:"y"`
	Z float32 `json:"z" db:"z"`
}

// TODO: rename to "Transform"
type ObjectTransform struct {
	Position Vec3 `db:"location" json:"location"`
	Rotation Vec3 `db:"rotation" json:"rotation"`
	Scale    Vec3 `db:"scale" json:"scale"`
}

type UserTransform struct {
	Position *Vec3 `db:"location" json:"location"`
	Rotation *Vec3 `db:"rotation" json:"rotation"`
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

func NewUserTransform() UserTransform {
	var t UserTransform
	t.Position = &Vec3{}
	t.Rotation = &Vec3{}
	return t
}

func (t *UserTransform) CopyToBuffer(b []byte) error {
	binary.LittleEndian.PutUint32(b, math.Float32bits(t.Position.X))
	binary.LittleEndian.PutUint32(b[Float32Bytes:], math.Float32bits(t.Position.Y))
	binary.LittleEndian.PutUint32(b[2*Float32Bytes:], math.Float32bits(t.Position.Z))
	binary.LittleEndian.PutUint32(b[3*Float32Bytes:], math.Float32bits(t.Rotation.X))
	binary.LittleEndian.PutUint32(b[4*Float32Bytes:], math.Float32bits(t.Rotation.Y))
	binary.LittleEndian.PutUint32(b[5*Float32Bytes:], math.Float32bits(t.Rotation.Z))
	return nil
}

func (t *UserTransform) Bytes() []byte {
	b := make([]byte, 6*Float32Bytes)
	t.CopyToBuffer(b)
	return b
}

func (t *UserTransform) CopyFromBuffer(b []byte) error {
	t.Position.X = math.Float32frombits(binary.LittleEndian.Uint32(b))
	t.Position.Y = math.Float32frombits(binary.LittleEndian.Uint32(b[Float32Bytes:]))
	t.Position.Z = math.Float32frombits(binary.LittleEndian.Uint32(b[2*Float32Bytes:]))
	t.Position.X = math.Float32frombits(binary.LittleEndian.Uint32(b[3*Float32Bytes:]))
	t.Position.Y = math.Float32frombits(binary.LittleEndian.Uint32(b[4*Float32Bytes:]))
	t.Position.Z = math.Float32frombits(binary.LittleEndian.Uint32(b[5*Float32Bytes:]))
	return nil
}

func (t *UserTransform) CopyTo(t1 *UserTransform) {
	*t1.Position = *t.Position
	*t1.Rotation = *t.Rotation
}
