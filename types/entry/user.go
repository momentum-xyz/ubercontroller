package entry

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID     uuid.UUID    `db:"user_id"`
	UserTypeID *uuid.UUID   `db:"user_type_id"`
	Profile    *UserProfile `db:"profile"`
	Options    *UserOptions `db:"options"`
	JWT        *JWT         `db:"auth"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  *time.Time   `db:"updated_at"`
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

type JWT struct {
	SignedString string    `db:"signed_string" json:"signed_string"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
}

type UserAttributeID struct {
	AttributeID
	UserID uuid.UUID `db:"user_id"`
}

type UserUserAttributeID struct {
	AttributeID
	SourceUserID uuid.UUID `db:"source_user_id"`
	TargetUserID uuid.UUID `db:"target_user_id"`
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

func NewUserAttributeID(attributeID AttributeID, userID uuid.UUID) UserAttributeID {
	return UserAttributeID{
		AttributeID: attributeID,
		UserID:      userID,
	}
}

func NewUserUserAttributeID(attributeID AttributeID, sourceUserID, targetUserID uuid.UUID) UserUserAttributeID {
	return UserUserAttributeID{
		AttributeID:  attributeID,
		SourceUserID: sourceUserID,
		TargetUserID: targetUserID,
	}
}
