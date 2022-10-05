package user

import (
	"github.com/influxdata/influxdb-client-go/v2/api/write"
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

	p.AddTag(prefix+"User", utils.AnonymizeUUID(u.id))

	return p

}
