package position_algo

import (
	"math"

	"github.com/momentum-xyz/controller/utils"
	cmath "github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

const (
	helixAngleDefaultValue       = 0.0
	helixSpiralScaleDefaultValue = 10.0
	helixRadiusDefaultValue      = 100.0
	helixHelixVShiftDefaultValue = 500.0
	helixVShiftDefaultValue      = 10.0
)

type helix struct {
	Angle       float64 `json:"angle"`
	SpiralScale float64 `json:"spiralScale"`
	R           float64 `json:"R"`
	VShift      float64 `json:"Vshift"`
	HelixVShift float64 `json:"helixVshift"`
}

func NewHelix(parameterMap map[string]interface{}) Algo {
	return &helix{
		Angle:       utils.GetFromAnyMap(parameterMap, "angle", helixAngleDefaultValue),
		VShift:      utils.GetFromAnyMap(parameterMap, "Vshift", helixVShiftDefaultValue),
		R:           utils.GetFromAnyMap(parameterMap, "R", helixRadiusDefaultValue),
		SpiralScale: utils.GetFromAnyMap(parameterMap, "spiralScale", helixSpiralScaleDefaultValue),
		HelixVShift: utils.GetFromAnyMap(parameterMap, "helixVshift", helixHelixVShiftDefaultValue),
	}
}

func (h *helix) CalcPos(parentTheta float64, parentPosition cmath.SpacePosition, i, n int) (
	cmath.SpacePosition, float64,
) {
	parent := parentPosition.Location.ToVec3f64()
	id := float64(i)

	acf := h.Angle / 360.0 * id
	r := h.R + h.SpiralScale*acf

	phi := 0.5*math.Pi + parentTheta
	angle := phi + id*h.Angle/180.0*math.Pi

	p := cmath.Vec3f64{
		X: math.Round((parent.X+r*math.Cos(angle))*10.0) / 10.0,
		Y: parent.Y + h.VShift + h.HelixVShift*acf,
		Z: math.Round((parent.Z+r*math.Sin(angle))*10.0) / 10.0,
	}

	np := cmath.SpacePosition{Location: p.ToVec3()}
	return np, math.Atan2(p.Z-parent.Z, p.X-parent.X) /* theta */
}

func (*helix) Name() string {
	return "helix"
}
