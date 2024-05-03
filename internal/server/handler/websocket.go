package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"password_keeper/internal/common/entity"
)

const timeOut = 10 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) handleSetWebsocket(w http.ResponseWriter, r *http.Request) {
	var data entity.WebDataMSG
	currUser := r.Context().Value("user").(int)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "can't upgrade to websocket", http.StatusInternalServerError)
		return
	}
	defer ws.Close()

	for {
		_, message, errMsg := ws.ReadMessage()
		if errMsg != nil {
			http.Error(w, "can't read message from websocket", http.StatusInternalServerError)
			break
		}

		err = json.Unmarshal(message, &data)
		if err != nil {
			http.Error(w, "can't parse message from websocket", http.StatusInternalServerError)
			break
		}

		switch data.Action {
		case "add":
			err = h.service.DataServiceInterface.SetData(currUser, data.Payload, data.EventType)
			if err != nil {
				http.Error(w, "err to set data", http.StatusInternalServerError)
				break
			}
		case "get":
			payloadDB, err := h.service.DataServiceInterface.GetData(currUser, data.EventType)
			if err != nil {
				http.Error(w, "err to get data", http.StatusInternalServerError)
				break
			}

			err = ws.WriteMessage(websocket.TextMessage, payloadDB)
			if err != nil {
				http.Error(w, "err to write data to websocket", http.StatusInternalServerError)
				break
			}
		case "delete":
			err = h.service.DataServiceInterface.DeleteData(currUser, data.EventType)
			if err != nil {
				http.Error(w, "err to delete data", http.StatusInternalServerError)
				break
			}
		}
	}
}
