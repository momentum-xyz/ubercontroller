package utils

import (
	"encoding/hex"
	"github.com/c0mm4nd/go-ripemd"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var salt umid.UMID
var nodeId umid.UMID

func SetAnonymizer(n umid.UMID, s umid.UMID) {
	nodeId = n
	salt = s
}

func AnonymizeUUID(userId umid.UMID) string {
	hash := ripemd.New128()
	for i := 0; i < 16; i += 2 {
		hash.Write(userId[i : i+2])
		hash.Write(nodeId[i : i+2])
		hash.Write(salt[i : i+2])
	}
	return hex.EncodeToString(hash.Sum(nil))
}
