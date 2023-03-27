package posbus

import "github.com/momentum-xyz/ubercontroller/logger"

type StringMapAny map[string]any

func (v StringMapAny) MarshalMUS(buf []byte) int {
	logger.L().Infof("*********** StringMapAny.MarshalMUS is not implemented yet!")
	return 0
}

func (v *StringMapAny) UnmarshalMUS(buf []byte) (int, error) {
	logger.L().Infof("*********** StringMapAny.UnmarshalMUS is not implemented yet!")
	return 0, nil
}

func (v StringMapAny) SizeMUS() int {
	logger.L().Infof("*********** StringMapAny.SizeMUS is not implemented yet!")
	return 0
}
