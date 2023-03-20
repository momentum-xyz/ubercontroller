//go:build ignore

package main

import (
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/ymz-ncnk/musgo/v2"
	"reflect"
)

func main() {
	musGo, err := musgo.New()

	if err != nil {
		panic(err)
	}
	unsafe := false // To generate safe code.
	err = musGo.Generate(reflect.TypeOf((*cmath.ObjectTransform)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*cmath.UserTransform)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*cmath.Vec3)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*cmath.Vec3f64)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
}