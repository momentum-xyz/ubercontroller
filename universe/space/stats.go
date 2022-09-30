package space

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/momentum-xyz/ubercontroller/universe"
	"time"
)

func (s *Space) AddSpaceTags(prefix string, p *write.Point) *write.Point {
	//Space UUID
	//Space Type
	//Space type UUID
	//Space name

	if prefix != "" {
		prefix += " "
	}
	p.AddTag(prefix+"Space UUID", s.id.String())
	p.AddTag(prefix+"Space Name", s.GetName())
	p.AddTag(prefix+"Space Type", s.spaceType.GetName())
	p.AddTag(prefix+"Space Type UUID", s.spaceType.GetID().String())
	return p

}

func (s *Space) sendSpaceEnterLeaveStats(user universe.User, value int) error {
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
