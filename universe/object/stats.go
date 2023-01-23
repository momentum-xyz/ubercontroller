package object

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/momentum-xyz/ubercontroller/universe"
	"time"
)

func (s *Object) AddObjectTags(prefix string, p *write.Point) *write.Point {
	if prefix != "" {
		prefix += " "
	}
	p.AddTag(prefix+"Object UUID", s.GetID().String())
	p.AddTag(prefix+"Object Name", s.GetName())
	p.AddTag(prefix+"Object Type", s.objectType.GetName())
	p.AddTag(prefix+"Object Type UUID", s.objectType.GetID().String())
	return p

}

func (s *Object) sendObjectEnterLeaveStats(user universe.User, value int) error {
	p := influxdb2.NewPoint(
		"object_join",
		map[string]string{},
		map[string]interface{}{"value": value},
		time.Now(),
	)
	user.AddInfluxTags("", p)
	s.AddObjectTags("", p)
	return s.world.WriteInfluxPoint(p)
}
