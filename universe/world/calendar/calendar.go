package calendar

import (
	"context"
	"time"

	"github.com/momentum-xyz/ubercontroller/types/entry"
)

type Calendar struct {
	ctx context.Context
}

func NewCalendar() *Calendar {
	calendar := &Calendar{}

	return calendar

}

func (c *Calendar) Initialize(ctx context.Context) error {
	c.ctx = ctx
	return nil
}

func (*Calendar) Run() error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//fmt.Println("Calendar timer ***")
		}
	}
	return nil
}

func (*Calendar) Stop() {

}

func (*Calendar) OnAttributeUpsert(attributeID entry.AttributeID, value *entry.AttributeValue) {

}

func (*Calendar) OnAttributeRemove(attributeID entry.AttributeID) {

}
