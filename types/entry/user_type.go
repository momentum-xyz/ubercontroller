package entry

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type UserType struct {
	UserTypeID   umid.UMID    `db:"user_type_id" json:"user_type_id"`
	UserTypeName string       `db:"user_type_name" json:"user_type_name"`
	Description  string       `db:"description" json:"description"`
	Options      *UserOptions `db:"options" json:"options"`
}
