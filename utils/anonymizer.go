package utils

import (
	"encoding/hex"
	"github.com/c0mm4nd/go-ripemd"
	"github.com/momentum-xyz/ubercontroller/utils/mid"
)

var salt mid.ID
var nodeId mid.ID

func SetAnonymizer(n mid.ID, s mid.ID) {
	nodeId = n
	salt = s
}

func AnonymizeUUID(userId mid.ID) string {
	hash := ripemd.New128()
	for i := 0; i < 16; i += 2 {
		hash.Write(userId[i : i+2])
		hash.Write(nodeId[i : i+2])
		hash.Write(salt[i : i+2])
	}
	return hex.EncodeToString(hash.Sum(nil))
}
