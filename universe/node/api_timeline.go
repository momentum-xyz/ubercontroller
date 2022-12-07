package node

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

const jsonStr = `[
  {
    "id": 1,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Jeroenski",
    "description": "",
    "image": "https://picsum.photos/100",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 3,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Polkajor",
    "description": "",
    "image": "https://picsum.photos/102",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "connected",
    "connectedTo": {
      "id": 4,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/104",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "connected"
    }
  },
  {
    "id": 2,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Brandskari",
    "description": "",
    "image": "https://picsum.photos/106",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 5,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Kidachu",
    "description": "",
    "image": "https://picsum.photos/108",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "docked",
    "dockedTo": {
      "id": 6,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/110",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "docked"
    }
  },
  {
    "id": 6,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Jeroenski",
    "description": "",
    "image": "https://picsum.photos/112",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 7,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Polkajor",
    "description": "",
    "image": "https://picsum.photos/114",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "connected",
    "connectedTo": {
      "id": 4,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/116",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "connected"
    }
  },
  {
    "id": 8,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Brandskari",
    "description": "",
    "image": "https://picsum.photos/118",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 9,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Kidachu",
    "description": "",
    "image": "https://picsum.photos/120",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "docked",
    "dockedTo": {
      "id": 6,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/122",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "docked"
    }
  },
  {
    "id": 10,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Jeroenski",
    "description": "",
    "image": "https://picsum.photos/124",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 11,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Polkajor",
    "description": "",
    "image": "https://picsum.photos/126",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "connected",
    "connectedTo": {
      "id": 4,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/128",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "connected"
    }
  },
  {
    "id": 12,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Brandskari",
    "description": "",
    "image": "https://picsum.photos/130",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 13,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Kidachu 2",
    "description": "",
    "image": "https://picsum.photos/132",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "docked",
    "dockedTo": {
      "id": 6,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/100",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "docked"
    }
  },
  {
    "id": 14,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Brandskari",
    "description": "",
    "image": "https://picsum.photos/102",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 15,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Kidachu 2",
    "description": "",
    "image": "https://picsum.photos/105",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "docked",
    "dockedTo": {
      "id": 6,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/107",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "docked"
    }
  }
]
    `

type UserItem struct {
	Id           int       `json:"id"`
	CollectionId int       `json:"collectionId"`
	Uuid         uuid.UUID `json:"uuid"`
	Owner        string    `json:"owner"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Image        string    `json:"image"`
	Date         time.Time `json:"date"`
	Type         string    `json:"type"`
}
type FeedItem struct {
	UserItem
	ConnectedTo   *UserItem  `json:"connectedTo,omitempty"`
	DockedTo      *UserItem  `json:"dockedTo,omitempty"`
	CalendarImage string     `json:"calendarImage"`
	CalendarTitle string     `json:"calendarTitle"`
	CalendarStart *time.Time `json:"calendarStart,omitempty"`
	CalendarEnd   *time.Time `json:"calendarEnd,omitempty"`
}

type Event struct {
	SpaceID     *uuid.UUID `json:"spaceId"`
	Title       string     `json:"title"`
	Start       time.Time  `json:"start"`
	End         time.Time  `json:"end"`
	EventID     string     `json:"eventId"`
	ImageHash   string     `json:"image"`
	WebLink     string     `json:"web_link"`
	Description string     `json:"description"`
	HostedBy    string     `json:"hosted_by"`
}

func (n *Node) apiNewsFeed(c *gin.Context) {

	list := make([]FeedItem, 0)

	err := json.Unmarshal([]byte(jsonStr), &list)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiNewsFeed: failed to Unmarshal sample jsonString")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
	}

	events, err := n.getAllEvents()
	if err != nil {
		err = errors.WithMessage(err, "Node: apiNewsFeed: failed to getAllEvents")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
	}

	for _, e := range events {

		userName := ""
		userAvatarHash := ""
		user, err := n.db.UsersGetUserByID(context.Background(), *e.SpaceID)
		if err != nil {
			err = errors.WithMessage(err, "Node: apiNewsFeed: failed to UsersGetUserByID")
			log.Warn(err)
		}

		if user != nil {
			if user.Profile != nil {
				if user.Profile.Name != nil {
					userName = *user.Profile.Name
				}

				if user.Profile.AvatarHash != nil {
					userAvatarHash = *user.Profile.AvatarHash
				}
			}
		}

		item := FeedItem{
			UserItem: UserItem{
				Id:           0,
				CollectionId: 0,
				Uuid:         *e.SpaceID,
				Owner:        "",
				Name:         userName, // User name userID=SpaceID
				Description:  e.Title,
				Image:        userAvatarHash, // User Avatar hash
				Date:         e.Start,
				Type:         "calendar_event"},
			ConnectedTo:   nil,
			DockedTo:      nil,
			CalendarTitle: e.Title,
			CalendarImage: e.ImageHash,
			CalendarStart: &e.Start,
			CalendarEnd:   &e.End,
		}
		list = append(list, item)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Date.Before(list[j].Date)
	})

	c.JSON(http.StatusOK, list)
}

func (n *Node) getAllEvents() ([]*Event, error) {
	attributes, err := n.db.SpaceAttributesGetSpaceAttributesByPluginIDAndAttributeName(context.Background(), universe.GetSystemPluginID(), "events")
	if err != nil {
		return nil, errors.WithMessage(err, "Node: getAllEvents: failed to SpaceAttributesGetSpaceAttributesByPluginIDAndAttributeName")
	}

	events := make([]*Event, 0)

	for _, attribute := range attributes {
		if attribute.Value != nil {
			for _, eventData := range *attribute.Value {
				e := getEvent(&attribute.SpaceID, eventData)
				if e != nil {
					events = append(events, e)
				}
			}
		}
	}

	return events, nil
}

func getEvent(spaceID *uuid.UUID, item any) *Event {
	e := &Event{SpaceID: spaceID}

	err := utils.MapDecode(item, e)
	if err != nil {
		log.Error(errors.WithMessage(err, "getEvent: failed to MapDecode 'events' attribute payload"))
		return nil
	}

	return e
}

func (n *Node) apiNotifications(c *gin.Context) {
	//TODO
	jsonStr := `[
  {
    "id": 1,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Jeroenski",
    "description": "",
    "image": "https://picsum.photos/100",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 3,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Polkajor",
    "description": "",
    "image": "https://picsum.photos/102",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "connected",
    "connectedTo": {
      "id": 4,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/104",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "connected"
    }
  },
  {
    "id": 2,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Brandskari",
    "description": "",
    "image": "https://picsum.photos/106",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 5,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Kidachu",
    "description": "",
    "image": "https://picsum.photos/108",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "docked",
    "dockedTo": {
      "id": 6,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/110",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "docked"
    }
  },
  {
    "id": 6,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Jeroenski",
    "description": "",
    "image": "https://picsum.photos/112",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 7,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Polkajor",
    "description": "",
    "image": "https://picsum.photos/114",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "connected",
    "connectedTo": {
      "id": 4,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/116",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "connected"
    }
  },
  {
    "id": 8,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Brandskari",
    "description": "",
    "image": "https://picsum.photos/118",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 9,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Kidachu",
    "description": "",
    "image": "https://picsum.photos/120",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "docked",
    "dockedTo": {
      "id": 6,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/122",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "docked"
    }
  },
  {
    "id": 10,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Jeroenski",
    "description": "",
    "image": "https://picsum.photos/124",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 11,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Polkajor",
    "description": "",
    "image": "https://picsum.photos/126",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "connected",
    "connectedTo": {
      "id": 4,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/128",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "connected"
    }
  },
  {
    "id": 12,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Brandskari",
    "description": "",
    "image": "https://picsum.photos/130",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 13,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Kidachu 2",
    "description": "",
    "image": "https://picsum.photos/132",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "docked",
    "dockedTo": {
      "id": 6,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/100",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "docked"
    }
  },
  {
    "id": 14,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Brandskari",
    "description": "",
    "image": "https://picsum.photos/102",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "created"
  },
  {
    "id": 15,
    "collectionId": 1,
    "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
    "owner": "",
    "name": "Kidachu 2",
    "description": "",
    "image": "https://picsum.photos/105",
    "date": "2022-11-25T08:05:48.447Z",
    "type": "docked",
    "dockedTo": {
      "id": 6,
      "collectionId": 1,
      "uuid": "d83670c7-a120-47a4-892d-f9ec75604f74",
      "owner": "",
      "name": "Space Odyssey",
      "description": "",
      "image": "https://picsum.photos/107",
      "date": "2022-11-25T08:05:48.447Z",
      "type": "docked"
    }
  }
]
    `
	c.DataFromReader(http.StatusOK,
		int64(len(jsonStr)), gin.MIMEJSON, strings.NewReader(jsonStr), nil)
}
