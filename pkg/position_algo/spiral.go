package position_algo

import (
	"math"

	"github.com/momentum-xyz/controller/utils"
	cmath "github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

const (
	spiralAngleDefaultValue  float64 = 0.0
	spiralRadiusDefaultValue float64 = 100.0
	spiralScaleDefaultValue  float64 = 10.0
	spiralVShiftDefaultValue float64 = 10.0
)

type spiral struct {
	Angle  float64 `json:"angle"`
	Scale  float64 `json:"spiralScale"`
	R      float64 `json:"R"`
	VShift float64 `json:"Vshift"`
}

func NewSpiral(parameterMap map[string]interface{}) Algo {
	return &spiral{
		Angle:  utils.GetFromAnyMap(parameterMap, "angle", spiralAngleDefaultValue),
		Scale:  utils.GetFromAnyMap(parameterMap, "spiralScale", spiralScaleDefaultValue),
		R:      utils.GetFromAnyMap(parameterMap, "R", spiralRadiusDefaultValue),
		VShift: utils.GetFromAnyMap(parameterMap, "Vshift", spiralVShiftDefaultValue),
	}
}

func (s *spiral) CalcPos(parentTheta float64, parentPosition cmath.SpacePosition, i, n int) (
	cmath.SpacePosition, float64,
) {
	parent := parentPosition.Location.ToVec3f64()
	scl := math.Sqrt(s.Scale * (float64(i) + s.Angle))

	r := s.R * scl

	phi := 0.5*math.Pi + parentTheta
	angle := phi + scl
	p := cmath.Vec3f64{
		X: math.Round((parent.X+r*math.Cos(angle))*10.0) / 10.0,
		Y: parent.Y + s.VShift,
		Z: math.Round((parent.Z+r*math.Sin(angle))*10.0) / 10.0,
	}

	np := cmath.SpacePosition{Location: p.ToVec3()}
	return np, math.Atan2(p.Z-parent.Z, p.X-parent.X) /* theta */
}

func (*spiral) Name() string {
	return "spiral"
}
