//go:build ignore

package main

import (
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/ymz-ncnk/musgo/v2"
	"reflect"
)

func main() {
	musGo, err := musgo.New()

	if err != nil {
		panic(err)
	}
	unsafe := false // To generate safe code.
	err = musGo.Generate(reflect.TypeOf((*dto.Asset3dType)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
}
