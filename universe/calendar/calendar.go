package calendar

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/posbus-protocol/posbus"

	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

type Calendar struct {
	ctx      context.Context
	world    universe.World
	timerSet *generic.TimerSet[string]
}

type Event struct {
	ObjectID *uuid.UUID `json:"objectId"`
	Title    string     `json:"title"`
	Start    time.Time  `json:"start"`
	End      time.Time  `json:"end"`
	EventID  string     `json:"eventId"`
}

var log = logger.L()

func NewCalendar(w universe.World) *Calendar {
	calendar := &Calendar{
		timerSet: generic.NewTimerSet[string](),
		world:    w,
	}

	return calendar
}

func (c *Calendar) Initialize(ctx context.Context) error {
	c.ctx = ctx
	return nil
}

func (c *Calendar) Run() error {
	go c.update()

	return nil
}

func (c *Calendar) update() {
	objects := c.world.GetAllObjects()

	events := getAllEvents(objects)
	nextEvents := findNextEvents(events)

	c.timerSet.StopAll()

	for i := range nextEvents {
		d := nextEvents[i].Start.Sub(time.Now())
		if d > 0 {
			c.timerSet.Set(nextEvents[i].EventID, d, c.tick)
		}
	}
}

func (c *Calendar) updateTimer() error {
	return nil
}

func (c *Calendar) tick(eventID string) error {
	e := c.getEventByID(eventID)
	if e == nil {
		return nil
	}
	topic := "notify-gathering-start"
	data, err := json.Marshal(&e)
	if err != nil {
		return errors.WithMessagef(err, "failed to marshal message payload")
	}
	m := posbus.NewRelayToReactMsg(topic, data).WebsocketMessage()
	c.world.Send(m, false)

	go c.update()

	return nil
}

func (c *Calendar) getEventByID(eventID string) *Event {
	objects := c.world.GetAllObjects()
	events := getAllEvents(objects)
	for _, e := range events {
		if e.EventID == eventID {
			return &e
		}
	}

	return nil
}

func findNextEvents(events []Event) []Event {
	if len(events) == 0 {
		return nil
	}

	// Filter out passed events
	result := make([]Event, 0)
	for _, e := range events {
		if e.Start.After(time.Now()) {
			result = append(result, e)
		}
	}

	if len(result) == 0 {
		return nil
	}

	min := result[0]
	for _, e := range result {
		if e.Start.Before(min.Start) {
			min = e
		}
	}

	// We can have several events starting at the same time
	result2 := make([]Event, 0)
	for _, e := range result {
		if e.Start.Equal(min.Start) {
			result2 = append(result2, e)
		}
	}

	return result2
}

func getAllEvents(objects map[uuid.UUID]universe.Object) []Event {
	attributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Object.Events.Name)

	//a := c.world.GetObjectAttributesValue(true)

	attributes := make([]*entry.AttributeValue, 0)
	events := make([]Event, 0)
	for objectID, _ := range objects {
		object := objects[objectID]

		attributeValue, ok := object.GetObjectAttributes().GetValue(attributeID)
		if !ok {
			continue
		}

		if attributeValue != nil {
			attributes = append(attributes, attributeValue)
			attribute := *attributeValue
			for _, event := range attribute {
				e, err := getEvent(&objectID, event)
				if err != nil {
					log.Error(err)
				}
				if e != nil {
					events = append(events, *e)
				}
			}
		}
	}

	return events
}

func getEvent(objectID *uuid.UUID, item any) (*Event, error) {
	e := &Event{ObjectID: objectID}

	err := utils.MapDecode(item, e)

	return e, errors.WithMessage(err, "utils.MapDecode")
}

func (*Calendar) Stop() error {
	return nil
}

func (c *Calendar) OnAttributeUpsert(attributeID entry.AttributeID, value any) {
	if attributeID.PluginID == universe.GetSystemPluginID() && attributeID.Name == universe.ReservedAttributes.Object.Events.Name {
		go c.update()
	}
}

func (*Calendar) OnAttributeRemove(attributeID entry.AttributeID) {

}
