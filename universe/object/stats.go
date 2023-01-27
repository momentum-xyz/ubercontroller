package object

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/momentum-xyz/ubercontroller/universe"
	"time"
)

func (o *Object) AddObjectTags(prefix string, p *write.Point) *write.Point {
	if prefix != "" {
		prefix += " "
	}
	p.AddTag(prefix+"Object UUID", o.GetID().String())
	p.AddTag(prefix+"Object Name", o.GetName())
	p.AddTag(prefix+"Object Type", o.objectType.GetName())
	p.AddTag(prefix+"Object Type UUID", o.objectType.GetID().String())
	return p

}

func (o *Object) sendObjectEnterLeaveStats(user universe.User, value int) error {
	p := influxdb2.NewPoint(
		"object_join",
		map[string]string{},
		map[string]interface{}{"value": value},
		time.Now(),
	)
	user.AddInfluxTags("", p)
	o.AddObjectTags("", p)
	return o.world.WriteInfluxPoint(p)
}
