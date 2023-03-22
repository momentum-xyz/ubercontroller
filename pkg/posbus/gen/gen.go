//go:build tools

package main

import (
	"bufio"
	"fmt"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/ymz-ncnk/musgo/v2"
	"os"
)

func main() {
	generateMus()
	generateTypes()
}

func generateMus() {
	musGo, err := musgo.New()

	if err != nil {
		panic(err)
	}
	unsafe := false // To generate safe code.

	for _, mId := range posbus.GetMessageIds() {
		posbus.MessageDataTypeById(mId)
		err = musGo.Generate(posbus.MessageDataTypeById(mId), unsafe)
		if err != nil {
			panic(err)
		}
	}

	for _, t := range posbus.ExtraTypes() {
		err = musGo.Generate(t, unsafe)
		if err != nil {
			panic(err)
		}
	}
}

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}

func generateTypes() {
	f, err := os.Create("types.autogen.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	_, err = fmt.Fprintf(w, "package posbus\n\nconst (\n")
	check_error(err)

	maxLen := 0
	for _, mId := range posbus.GetMessageIds() {
		l := len(posbus.MessageTypeNameById(mId))
		if l > maxLen {
			maxLen = l
		}
	}

	for _, mId := range posbus.GetMessageIds() {
		mTypeName := posbus.MessageTypeNameById(mId)
		_, err = fmt.Fprintf(w, "\t%-*sMsgType = 0x%08X\n", maxLen+5, "Type"+mTypeName, mId)
		check_error(err)
	}
	_, err = fmt.Fprintf(w, ")\n")
	check_error(err)
	w.Flush()
}
