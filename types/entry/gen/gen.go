//go:build tools

package main

//
//import (
//	"github.com/momentum-xyz/ubercontroller/types/entry"
//	"github.com/ymz-ncnk/musgo/v2"
//	"reflect"
//)
//
//func main() {
//	musGo, err := musgo.New()
//
//	if err != nil {
//		panic(err)
//	}
//	unsafe := false // To generate safe code.
//
//	err = musGo.Generate(reflect.TypeOf((*entry.AttributeValue)(nil)).Elem(), unsafe)
//	if err != nil {
//		panic(err)
//	}
//}
