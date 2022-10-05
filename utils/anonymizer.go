package utils

import (
	"encoding/hex"
	"github.com/c0mm4nd/go-ripemd"
	"github.com/google/uuid"
)

var salt uuid.UUID
var nodeId uuid.UUID

func SetAnonymizer(n uuid.UUID, s uuid.UUID) {
	nodeId = n
	salt = s
}

func AnonymizeUUID(userId uuid.UUID) string {
	hash := ripemd.New128()
	for i := 0; i < 16; i += 2 {
		hash.Write(userId[i : i+2])
		hash.Write(nodeId[i : i+2])
		hash.Write(salt[i : i+2])
	}
	return hex.EncodeToString(hash.Sum(nil))
}
