package user

import (
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func (u *User) AddInfluxTags(prefix string, p *write.Point) *write.Point {
	userTypeId := u.GetUserType().GetID()
	userType := u.GetUserType().GetName()

	if prefix != "" {
		prefix += " "
	}
	p.AddTag(prefix+"User Type UUID", userTypeId.String())
	p.AddTag(prefix+"User Type", userType)

	p.AddTag(prefix+"User", utils.AnonymizeUUID(u.GetID()))

	return p
}

func (u *User) SendHighFiveStats(target universe.User) error {
	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		if current == nil {
			current = entry.NewAttributePayload(nil, nil)
		}
		if current.Value == nil {
			current.Value = entry.NewAttributeValue()
		}

		// increment value of high five counter by 1
		(*current.Value)[universe.Attributes.User.HighFive.Key] = utils.GetFromAnyMap(
			*current.Value, universe.Attributes.User.HighFive.Key, float64(0),
		) + 1

		return current, nil
	}

	if _, err := universe.GetNode().UpsertUserUserAttribute(
		entry.NewUserUserAttributeID(
			entry.NewAttributeID(
				universe.GetSystemPluginID(), universe.Attributes.User.HighFive.Name,
			),
			u.GetID(), target.GetID(),
		), modifyFn,
	); err != nil {
		return errors.New("failed to upsert high-five user user attribute")
	}

	return nil
}
