package node

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	Uuid         string    `json:"uuid"`
	Owner        string    `json:"owner"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Image        string    `json:"image"`
	Date         time.Time `json:"date"`
	Type         string    `json:"type"`
}
type FeedItem struct {
	UserItem
	ConnectedTo *UserItem `json:"connectedTo,omitempty"`
	DockedTo    *UserItem `json:"dockedTo,omitempty"`
}

func (n *Node) apiNewsFeed(c *gin.Context) {

	list := make([]FeedItem, 0)

	json.Unmarshal([]byte(jsonStr), &list)

	//TODO

	c.JSON(http.StatusOK, list)
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
