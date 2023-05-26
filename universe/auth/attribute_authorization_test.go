package auth

import (
	"context"
	"reflect"
	"testing"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

func TestCheckAttributePermissions(t *testing.T) {
	pluginID := umid.MustParse("00000000-0000-8000-8000-000000000001")
	userID := umid.MustParse("00000000-0000-8000-8000-000000000002")
	//objectId := umid.MustParse("00000000-0000-8000-8000-000000000003")
	attrTypeID := entry.AttributeTypeID{
		PluginID: pluginID,
		Name:     "test-attribute",
	}
	attrID := entry.AttributeID{ // e.g. 'ID' for ObjectAttribute, a bit weird.
		PluginID: pluginID,
		Name:     "test-attribute",
	}
	attrTypeDesc := "Attribute for testing"

	tests := []struct {
		name        string
		permissions *entry.PermissionsAttributeOption
		opType      operationType
		userRoles   []entry.PermissionsRoleType
		want        bool
		wantErr     bool
	}{
		{
			name:        "empty permissions",
			permissions: nil,
			opType:      ReadOperation,
			want:        false,
			wantErr:     false,
		},
		{
			name: "admin user read on any permission",
			permissions: &entry.PermissionsAttributeOption{
				Read:  "any",
				Write: "admin+user_owner",
			},
			userRoles: []entry.PermissionsRoleType{
				entry.PermissionAdmin,
			},
			opType:  ReadOperation,
			want:    true,
			wantErr: false,
		},
		{
			name: "admin user write",
			permissions: &entry.PermissionsAttributeOption{
				Read:  "any",
				Write: "admin+user_owner",
			},
			userRoles: []entry.PermissionsRoleType{
				entry.PermissionAdmin,
			},
			opType:  WriteOperation,
			want:    true,
			wantErr: false,
		},
		{
			name: "guest user write",
			permissions: &entry.PermissionsAttributeOption{
				Read:  "any",
				Write: "admin+user_owner",
			},
			userRoles: []entry.PermissionsRoleType{
				entry.PermissionUser, // currently now specific guest role
			},
			opType:  WriteOperation,
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			aType := entry.AttributeType{
				AttributeTypeID: attrTypeID,
				Description:     &attrTypeDesc,
				Options: &entry.AttributeOptions{
					"permissions": tt.permissions,
				},
			}
			attrStore := getAttrStoreStub(aType, tt.userRoles)

			got, err := CheckAttributePermissions(
				ctx,
				aType, attrStore, attrID,
				userID, tt.opType)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAttributePermissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAttributePermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

type attrStoreStub struct {
	AT        entry.AttributeType
	UserRoles []entry.PermissionsRoleType
}

func (a attrStoreStub) GetOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	return a.AT.Options, true
}

func (a attrStoreStub) GetEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	return a.GetOptions(attributeID)
}

func (a attrStoreStub) GetUserRoles(
	ctx context.Context,
	attrType entry.AttributeType,
	attrID entry.AttributeID,
	userID umid.UMID,
) ([]entry.PermissionsRoleType, error) {
	return a.UserRoles, nil
}

func getAttrStoreStub(at entry.AttributeType, roles []entry.PermissionsRoleType) AttributePermissionsAuthorizer[entry.AttributeID] {
	return attrStoreStub{at, roles}
}

func TestPermissionsFromOptions(t *testing.T) {
	tests := []struct {
		name    string
		options entry.AttributeOptions
		want    *entry.PermissionsAttributeOption
		wantErr bool
	}{
		{
			options: nil,
			want:    nil,
		},
		{
			options: entry.AttributeOptions{},
			want:    nil,
		},
		{
			options: entry.AttributeOptions{
				"render_auto": map[string]any{"foo": "bar"},
				"permissions": map[string]any{
					"read":  "any",
					"write": "admin+user_owner",
				},
			},
			want: &entry.PermissionsAttributeOption{
				Read: "any", Write: "admin+user_owner",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := permissionsFromOptions(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("PermissionsFromOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PermissionsFromOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
