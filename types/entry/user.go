package entry

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"time"
)

type User struct {
	UserID     umid.UMID      `db:"user_id" json:"user_id"`
	UserTypeID umid.UMID      `db:"user_type_id" json:"user_type_id"`
	Profile    UserProfile    `db:"profile" json:"profile"`
	Options    *UserOptions   `db:"options" json:"options"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at" json:"updated_at"`
	Auth       map[string]any `db:"auth" json:"auth"`
}

type UserOptions struct {
	IsGuest *bool `db:"is_guest" json:"is_guest"`
}

type UserProfile struct {
	Name        *string `db:"name" json:"name"`
	Bio         *string `db:"bio" json:"bio"`
	Location    *string `db:"location" json:"location"`
	AvatarHash  *string `db:"avatar_hash" json:"avatar_hash"`
	ProfileLink *string `db:"profile_link" json:"profile_link"`
	OnBoarded   *bool   `db:"onboarded" json:"onboarded"`
}

type UserAttributeID struct {
	AttributeID
	UserID umid.UMID `db:"user_id" json:"user_id"`
}

type UserUserAttributeID struct {
	AttributeID
	SourceUserID umid.UMID `db:"source_user_id" json:"source_user_id"`
	TargetUserID umid.UMID `db:"target_user_id" json:"target_user_id"`
}

type UserAttribute struct {
	UserAttributeID
	*AttributePayload
}

type UserUserAttribute struct {
	UserUserAttributeID
	*AttributePayload
}

func NewUserAttribute(userAttributeID UserAttributeID, payload *AttributePayload) *UserAttribute {
	return &UserAttribute{
		UserAttributeID:  userAttributeID,
		AttributePayload: payload,
	}
}

func NewUserUserAttribute(userUserAttributeID UserUserAttributeID, payload *AttributePayload) *UserUserAttribute {
	return &UserUserAttribute{
		UserUserAttributeID: userUserAttributeID,
		AttributePayload:    payload,
	}
}

func NewUserAttributeID(attributeID AttributeID, userID umid.UMID) UserAttributeID {
	return UserAttributeID{
		AttributeID: attributeID,
		UserID:      userID,
	}
}

func NewUserUserAttributeID(attributeID AttributeID, sourceUserID, targetUserID umid.UMID) UserUserAttributeID {
	return UserUserAttributeID{
		AttributeID:  attributeID,
		SourceUserID: sourceUserID,
		TargetUserID: targetUserID,
	}
}
