//go:generate go run gen/mus.go

package umid

import (
	"database/sql/driver"
	"encoding/binary"
	"github.com/gofrs/uuid/v5"
	_ "github.com/ymz-ncnk/muserrs"
	_ "github.com/ymz-ncnk/musgo/v2"
)

var Nil = UMID{}

type UMID uuid.UUID

func Parse(s string) (UMID, error) {
	//r, e := uuid.Parse(s)
	r, e := uuid.FromString(s)
	return UMID(r), e
}

func ParseBytes(b []byte) (UMID, error) {
	//r, e := uuid.ParseBytes(b)
	r, e := uuid.FromBytes(b)
	return UMID(r), e
}

func MustParse(s string) UMID {
	//return UMID(uuid.MustParse(s))
	return UMID(uuid.Must(uuid.FromString(s)))

}

func (id UMID) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *UMID) UnmarshalText(data []byte) error {
	return (*uuid.UUID)(id).UnmarshalText(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (id UMID) MarshalBinary() ([]byte, error) {
	return uuid.UUID(id).MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (id *UMID) UnmarshalBinary(data []byte) error {
	return (*uuid.UUID)(id).UnmarshalBinary(data)
}

func FromBytes(b []byte) (id UMID, err error) {
	r, e := uuid.FromBytes(b)
	return UMID(r), e
}

func Must(id UMID, err error) UMID {
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (id UMID) String() string {
	return uuid.UUID(id).String()
}

//// URN returns the RFC 2141 URN form of uuid,
//// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx,  or "" if uuid is invalid.
//func (umid UMID) URN() string {
//	return uuid.UMID(umid).URN()
//}

func (id UMID) Variant() byte {
	return uuid.UUID(id).Variant()
}

// Version returns the version of uuid.
func (id UMID) Version() byte {
	return uuid.UUID(id).Version()
}

func (id *UMID) Scan(src interface{}) error {
	return (*uuid.UUID)(id).Scan(src)
}

func (id UMID) Value() (driver.Value, error) {
	return uuid.UUID(id).Value()
}

func New() UMID {
	r, _ := uuid.NewV7()
	return UMID(r)
}

func (id UMID) ClockSequence() int {
	//	return uuid.UUID(umid).ClockSequence()
	return int(binary.BigEndian.Uint16(id[8:10])) & 0x3fff
}
