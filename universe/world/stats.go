package world

import (
	influx_write "github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func (w *World) WriteInfluxPoint(point *influx_write.Point) error {
	point.AddTag("World ID", w.GetID().String())
	point.AddTag("World Name", w.GetName())
	return universe.GetNode().WriteInfluxPoint(point)
}
