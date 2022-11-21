package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
)

type Calendar struct {
	ctx        context.Context
	world      universe.World
	nextEvents []Event
	timerSet   *generic.TimerSet[string]
}

type Event struct {
	SpaceID *uuid.UUID
	Title   string
	Start   time.Time
	End     time.Time
	EventID string
}

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

	c.nextEvents = nextEvents

	for _, e := range c.nextEvents {
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
	pluginID := uuid.MustParse("f0f0f0f0-0f0f-4ff0-af0f-f0f0f0f0f0f0")
	attributeID := entry.NewAttributeID(pluginID, "events")

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
				e := getEvent(&spaceID, event)
				if e != nil {
					events = append(events, *e)
				}
			}
		}
	}

	return events
}

func getEvent(spaceID *uuid.UUID, item any) *Event {
	e := &Event{SpaceID: spaceID}

	i, ok := item.(map[string]any)
	if !ok {
		return nil
	}
	start, ok := i["start"].(string)
	if ok {
		layout := "2006-01-02T15:04:05Z0700"
		t, err := time.Parse(layout, start)
		if err == nil {
			e.Start = t
		}
	}

	title, ok := i["title"].(string)
	if ok {
		e.Title = title
	}

	eventID, ok := i["eventId"].(string)
	if ok {
		e.EventID = eventID
	}

	return e
}

func (*Calendar) Stop() error {
	return nil
}

func (*Calendar) OnAttributeUpsert(attributeID entry.AttributeID, value *entry.AttributeValue) {

}

func (*Calendar) OnAttributeRemove(attributeID entry.AttributeID) {

}
