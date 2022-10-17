package position_algo

import (
	"math"

	"github.com/momentum-xyz/controller/utils"
	cm "github.com/momentum-xyz/ubercontroller/pkg/cmath"
)

const (
	sectorAngleDefaultValue  = 0.0
	sectorRadiusDefaultValue = 100.0
	sectorVShiftDefaultValue = 10.0
)

type sector struct {
	Angle  float64 `json:"angle"`
	R      float64 `json:"R"`
	VShift float64 `json:"Vshift"`
}

func NewSector(parameterMap map[string]interface{}) Algo {
	return &sector{
		Angle:  utils.GetFromAnyMap(parameterMap, "angle", sectorAngleDefaultValue),
		R:      utils.GetFromAnyMap(parameterMap, "R", sectorRadiusDefaultValue),
		VShift: utils.GetFromAnyMap(parameterMap, "Vshift", sectorVShiftDefaultValue),
	}
}

func (sec *sector) CalcPos(parentTheta float64, parentVector cm.Vec3, i, n int) (cm.Vec3, float64) {
	parent := parentVector.ToVec3f64()

	scl := float64(0)
	phi := parentTheta
	if n > 1 {
		scl += 2.0 * sec.Angle / 180.0 * math.Pi / float64(n-1)
		phi += -sec.Angle / 180.0 * math.Pi
	}

	angle := phi + float64(i)*scl
	p := cm.Vec3f64{
		X: math.Round((parent.X+sec.R*math.Cos(angle))*10.0) / 10.0,
		Y: parent.Y + sec.VShift,
		Z: math.Round((parent.Z+sec.R*math.Sin(angle))*10.0) / 10.0,
	}

	return p.ToVec3(), math.Atan2(p.Z-parent.Z, p.X-parent.X) /* theta */
}

func (*sector) Name() string {
	return "sector"
}
