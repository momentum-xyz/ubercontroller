//go:build tools

package main

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/ymz-ncnk/musgo/v2"
	"reflect"
)

func main() {
	musGo, err := musgo.New()

	if err != nil {
		panic(err)
	}
	unsafe := false // To generate safe code.
	//var id qqq.IDType
	////// Alias types don't support tags, so to set up validator we use
	////// GenerateAliasAs() method.
	//conf := musgo.DefAliasConf
	//conf.MaxLength = 16 // Restricts length of StringAlias values to 5 characters.
	//err = musGo.GenerateAliasAs(reflect.TypeOf(id), conf)
	//if err != nil {
	//	panic(err)
	//}

	err = musGo.Generate(reflect.TypeOf((*umid.UMID)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
}
