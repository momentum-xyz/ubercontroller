package node

import (
	"context"
	influxWrite "github.com/influxdata/influxdb-client-go/v2/api/write"
	"time"
)

const (
	influxDBTimeout = 2 * time.Second
)

func (n *Node) WriteInfluxPoint(point *influxWrite.Point) error {
	// TODO: uncomment once influx part ready
	// TODO: pass stats as simple map and unwrap here
	return nil
	point.AddTag("Node ID", n.GetID().String())
	point.AddTag("Node Name", n.GetName())
	ctx, cancel := context.WithTimeout(context.Background(), influxDBTimeout)
	defer func() {
		cancel()
		n.log.Warn("ControllerHub: WriteInfluxPoint: stat sent")
	}()

	return n.influx.WritePoint(ctx, point)
}
