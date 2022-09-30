package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, httpStatus int, response any) {
	// response must be pointer
	data, err := json.Marshal(response)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(httpStatus)
	w.Write(data)
}

func JSONError(w http.ResponseWriter, httpStatus int, err error) {
	w.WriteHeader(httpStatus)
	w.Write([]byte(fmt.Sprintf("{\"error\": %q}", err.Error())))
}

var WebsocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
