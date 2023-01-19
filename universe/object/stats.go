package object

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/momentum-xyz/ubercontroller/universe"
	"time"
)

func (s *Object) AddSpaceTags(prefix string, p *write.Point) *write.Point {
	if prefix != "" {
		prefix += " "
	}
	p.AddTag(prefix+"Object UUID", s.GetID().String())
	p.AddTag(prefix+"Object Name", s.GetName())
	p.AddTag(prefix+"Object Type", s.spaceType.GetName())
	p.AddTag(prefix+"Object Type UUID", s.spaceType.GetID().String())
	return p

}

func (s *Object) sendSpaceEnterLeaveStats(user universe.User, value int) error {
	p := influxdb2.NewPoint(
		"space_join",
		map[string]string{},
		map[string]interface{}{"value": value},
		time.Now(),
	)
	user.AddInfluxTags("", p)
	s.AddSpaceTags("", p)
	return s.world.WriteInfluxPoint(p)
}
