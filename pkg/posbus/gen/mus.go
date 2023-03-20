//go:build tools

package main

import (
	"bufio"
	"fmt"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/ymz-ncnk/musgo/v2"
	"os"
	"reflect"
	"sort"
)

func main() {
	musGo, err := musgo.New()

	if err != nil {
		panic(err)
	}
	unsafe := false // To generate safe code.
	err = musGo.Generate(reflect.TypeOf((*posbus.HandShake)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
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
	err = musGo.Generate(reflect.TypeOf((*posbus.UserData)(nil)).Elem(), unsafe)
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
	err = musGo.Generate(reflect.TypeOf((*posbus.SetWorld)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	//err = musGo.Generate(reflect.TypeOf((*posbus.ObjectData)(nil)).Elem(), unsafe)
	//if err != nil {
	//	panic(err)
	//}
	err = musGo.Generate(reflect.TypeOf((*posbus.LockObject)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.ObjectLockResult)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.ObjectPosition)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.UsersTransformList)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.UserTransform)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.MyTransform)(nil)).Elem(), unsafe)
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
	err = musGo.Generate(reflect.TypeOf((*posbus.GenericMessage)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.NotificationType)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}
	err = musGo.Generate(reflect.TypeOf((*posbus.Notification)(nil)).Elem(), unsafe)
	if err != nil {
		panic(err)
	}

	printTypes()
}

func printTypes() {
	f, err := os.Create("types.autogen.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	_, err = fmt.Fprintf(w, "package posbus\n\nconst (\n")
	if err != nil {
		panic(err)
	}

	msgNames := make([]string, 0, len(posbus.IdsCheck))
	typeMap := make(map[string]posbus.MsgType)

	for id, name := range posbus.IdsCheck {
		msgNames = append(msgNames, name)
		typeMap[name] = id
	}

	sort.Strings(msgNames)

	for _, msgName := range msgNames {
		msgType := typeMap[msgName]
		_, err = fmt.Fprintf(w, "    %+v     MsgType = 0x%08X\n", msgName, msgType)
		if err != nil {
			panic(err)
		}
	}
	_, err = fmt.Fprintf(w, ")\n")
	if err != nil {
		panic(err)
	}
	w.Flush()
}
