package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/pkg/errors"

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
	SpaceID *uuid.UUID
	Title   string
	Start   time.Time
	End     time.Time
	EventID string
}

var pluginID = uuid.MustParse("f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0")
var attributeName = "events"
var log = logger.L()

func NewCalendar() *Calendar {
	calendar := &Calendar{
		timerSet: generic.NewTimerSet[string](),
	}

	return calendar
}

func (c *Calendar) Initialize(ctx context.Context, w universe.World) error {
	c.ctx = ctx
	c.world = w
	return nil
}

func (c *Calendar) Run() error {
	fmt.Println("RUN calendar" + c.world.GetID().String())

	c.update()

	return nil
}

func (c *Calendar) update() {
	spaces := c.world.GetSpaces(true)
	events := getAllEvents(spaces)
	nextEvents := findNextEvents(events)

	c.timerSet.StopAll()

	for _, e := range nextEvents {
		d := e.Start.Sub(time.Now())
		if d > 0 {
			c.timerSet.Set(e.EventID, d, c.tick)
		}
	}
}

func (c *Calendar) updateTimer() error {
	return nil
}

func (c *Calendar) tick(eventID string) error {
	fmt.Println("TICK", eventID)

	e := c.getEventByID(eventID)
	topic := "notify-gathering-start"
	data, err := json.Marshal(&e)
	if err != nil {
		return errors.WithMessagef(err, "failed to marshal message payload")
	}
	m := posbus.NewRelayToReactMsg(topic, data).WebsocketMessage()
	c.world.Send(m, false)

	c.update()
	return nil
}

func (c *Calendar) getEventByID(eventID string) *Event {
	spaces := c.world.GetSpaces(true)
	events := getAllEvents(spaces)
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

	result2 := make([]Event, 0)
	for _, e := range result {
		if e.Start.Equal(min.Start) {
			result2 = append(result2, e)
		}
	}

	return result2
}

func getAllEvents(spaces map[uuid.UUID]universe.Space) []Event {
	attributeID := entry.NewAttributeID(pluginID, attributeName)

	//a := c.world.GetSpaceAttributesValue(true)

	attributes := make([]*entry.AttributeValue, 0)
	events := make([]Event, 0)
	for spaceID, _ := range spaces {
		space := spaces[spaceID]

		attributeValue, ok := space.GetSpaceAttributeValue(attributeID)
		if !ok {
			continue
		}

		if attributeValue != nil {
			attributes = append(attributes, attributeValue)
			attribute := *attributeValue
			for _, event := range attribute {
				e, err := getEvent(&spaceID, event)
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

func getEvent(spaceID *uuid.UUID, item any) (*Event, error) {
	e := &Event{SpaceID: spaceID}

	err := utils.MapDecode(item, e)

	return e, errors.WithMessage(err, "utils.MapDecode")
}

func (*Calendar) Stop() error {
	return nil
}

func (c *Calendar) OnAttributeUpsert(attributeID entry.AttributeID, value *entry.AttributeValue) {
	fmt.Println("OnAttributeUpsert ***", attributeID)
	if attributeID.PluginID == pluginID && attributeID.Name == attributeName {
		go c.update()
	}
}

func (*Calendar) OnAttributeRemove(attributeID entry.AttributeID) {

}
