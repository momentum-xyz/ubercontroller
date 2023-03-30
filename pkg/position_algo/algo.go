package position_algo

import "github.com/momentum-xyz/ubercontroller/pkg/cmath"

const (
	defaultVShiftValue = 10.0
	defaultAngleValue  = 0.0
	defaultRadiusValue = 100.0
	defaultSpiralScale = 10.0
	defaultHelixVShift = 500.0
)

type Algo interface {
	Name() string
	CalcPos(parentTheta float64, parentPosition cmath.Transform, i, n int) (cmath.Transform, float64)
}
