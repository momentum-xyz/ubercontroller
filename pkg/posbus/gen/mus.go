//go:build ignore

package main

import (
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/ymz-ncnk/musgo/v2"
	"reflect"
)

func main() {
	musGo, err := musgo.New()

	if err != nil {
		panic(err)
	}
	unsafe := false // To generate safe code.
	err = musGo.Generate(reflect.TypeOf((*posbus.AddObjects)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.RemoveObjects)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.AddUsers)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.UserDefinition)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.RemoveUsers)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.Signal)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.TeleportRequest)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.SetWorldData)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	//err = musGo.Generate(reflect.TypeOf((*posbus.ObjectData)(nil)).Elem(), unsafe)
	//if err != nil {
	//	panic(err)
	//}
	err = musGo.Generate(reflect.TypeOf((*posbus.SetObjectLock)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.ObjectLockResultData)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.ObjectPosition)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.SetUsersTransforms)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.UserTransform)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.SignalType)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.ObjectDefinition)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.GenericMessageData)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
}
