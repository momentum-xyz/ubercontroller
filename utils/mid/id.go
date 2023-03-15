package mid

import (
	"database/sql/driver"
	"encoding/binary"
	"github.com/gofrs/uuid/v5"
)

var Nil = ID{}

type ID uuid.UUID

func Parse(s string) (ID, error) {
	//r, e := uuid.Parse(s)
	r, e := uuid.FromString(s)
	return ID(r), e
}

func ParseBytes(b []byte) (ID, error) {
	//r, e := uuid.ParseBytes(b)
	r, e := uuid.FromBytes(b)
	return ID(r), e
}

func MustParse(s string) ID {
	//return ID(uuid.MustParse(s))
	return ID(uuid.Must(uuid.FromString(s)))

}

func (id ID) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *ID) UnmarshalText(data []byte) error {
	return (*uuid.UUID)(id).UnmarshalText(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (id ID) MarshalBinary() ([]byte, error) {
	return uuid.UUID(id).MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (id *ID) UnmarshalBinary(data []byte) error {
	return (*uuid.UUID)(id).UnmarshalBinary(data)
}

func FromBytes(b []byte) (id ID, err error) {
	r, e := uuid.FromBytes(b)
	return ID(r), e
}

func Must(id ID, err error) ID {
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (id ID) String() string {
	return uuid.UUID(id).String()
}

//// URN returns the RFC 2141 URN form of uuid,
//// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx,  or "" if uuid is invalid.
//func (mid ID) URN() string {
//	return uuid.ID(mid).URN()
//}

func (id ID) Variant() byte {
	return uuid.UUID(id).Variant()
}

// Version returns the version of uuid.
func (id ID) Version() byte {
	return uuid.UUID(id).Version()
}

func (id *ID) Scan(src interface{}) error {
	return (*uuid.UUID)(id).Scan(src)
}

func (id ID) Value() (driver.Value, error) {
	return uuid.UUID(id).Value()
}

func New() ID {
	r, _ := uuid.NewV7()
	return ID(r)
}

func (id ID) ClockSequence() int {
	//	return uuid.UUID(mid).ClockSequence()
	return int(binary.BigEndian.Uint16(id[8:10])) & 0x3fff
}
