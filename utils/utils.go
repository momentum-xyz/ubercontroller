package utils

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func BinID(id uuid.UUID) []byte {
	binID, err := id.MarshalBinary()
	if err != nil {
		log.Errorf("Utils: BinID: failed to marshal binary: %+v", errors.WithStack(err))
		return nil
	}
	return binID
}
