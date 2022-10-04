package node

import (
	"context"
	"encoding/hex"
	"github.com/c0mm4nd/go-ripemd"
	"github.com/google/uuid"
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

func (n *Node) HashID(userId uuid.UUID) string {
	hash := ripemd.New128()
	nodeId := n.GetID()
	for i := 0; i < 16; i += 2 {
		hash.Write(userId[i : i+2])
		hash.Write(nodeId[i : i+2])
		// TODO: fix once salt is read in attributes
		//hash.Write(n.userHashSalt[i : i+2])
	}
	return hex.EncodeToString(hash.Sum(nil))
}
