package utils

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var MASK_V8 = [16]byte{0, 0, 0, 0, 0, 0, 128, 0, 128, 0, 0, 0, 0, 0, 0, 0}

func BinID(id umid.UMID) []byte {
	binID, err := id.MarshalBinary()
	if err != nil {
		// TODO: return the error!
		log.Printf("Utils: BinID: failed to marshal binary: %+v", errors.WithStack(err))
		return nil
	}
	return binID
}

func MergeMaps[K comparable, V any](m1, m2 map[K]V) map[K]V {
	m := make(map[K]V, len(m1)+len(m2))
	for k, v := range m1 {
		m[k] = v
	}
	for k, v := range m2 {
		m[k] = v
	}
	return m
}

func MapDecode(input, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName: "json",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			handleNilAnonymousNestedStruct(),
			stringToUUIDHookFunc(),
			stringToTimeHookFunc(),
			mapToStringHookFunc(),
		),
		Squash: true,
		Result: output,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func MapEncode(input, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName: "json",
		//DecodeHook: mapstructure.ComposeDecodeHookFunc(
		//	handleNilAnonymousNestedStruct(),
		//	stringToUUIDHookFunc(),
		//	stringToTimeHookFunc(),
		//	mapToStringHookFunc(),
		//),
		Squash: true,
		Result: output,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func HexToAddress(s string) []byte {
	b, err := hex.DecodeString(s[2:])
	if err != nil {
		panic(err)
	}
	return b
}

func AddressToHex(a []byte) string {
	return hex.EncodeToString(a)
}

// handleNilAnonymousNestedStruct needed to fix "unsupported type for squash: ptr" mapstructure error
func handleNilAnonymousNestedStruct() mapstructure.DecodeHookFunc {
	return func(from reflect.Value, to reflect.Value) (interface{}, error) {
		if to.Kind() != reflect.Struct {
			return from.Interface(), nil
		}

		for i := 0; i < to.NumField(); i++ {
			fieldVal := to.Field(i)
			if !fieldVal.IsValid() {
				continue
			}
			if fieldVal.Kind() != reflect.Ptr {
				continue
			}
			if !fieldVal.IsNil() {
				continue
			}
			if !to.Type().Field(i).Anonymous {
				continue
			}

			fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		}

		return from.Interface(), nil
	}
}

func stringToUUIDHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(umid.UMID{}) {
			return data, nil
		}

		return umid.Parse(data.(string))
	}
}

func mapToStringHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}
		if t.Kind() != reflect.String {
			return data, nil
		}

		bytes, err := json.Marshal(data)
		return string(bytes), err
	}
}

func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		layout := "2006-01-02T15:04:05Z0700"

		return time.Parse(layout, data.(string))
	}
}

func GoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine umid: %v", err))
	}
	return id
}

func UMIDToSEQ(id umid.UMID) uint64 {
	var buf [16]byte
	for i := 0; i < 16; i++ {
		buf[i] = id[i] &^ MASK_V8[i]
	}
	return binary.BigEndian.Uint64(buf[8:])
}

func SEQtoUMID(s uint64) umid.UMID {
	id := umid.Nil
	var seq [16]byte
	binary.BigEndian.PutUint64(seq[8:], s)
	for i := 0; i < 16; i++ {
		id[i] = seq[i] | MASK_V8[i]
	}
	return id
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func AnyContains[T comparable](s []any, e T) bool {
	for _, a := range s {
		if a.(T) == e {
			return true
		}
	}
	return false
}
