package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}

func HandleSocketPayloadEvents(client *Client, socketEvent SocketEvent) {
	var socketEventResponse SocketEvent
	selectedUserID := socketEvent.EventPayload.UserID
	socketEventResponse.EventName = "Message Response"
	socketEventResponse.EventPayload = Payload{
		UserID: selectedUserID,
		Message: socketEvent.EventPayload.Message,
		Username: ,
	}
}

func EmitToSpecificClient(pool *Pool, payload SocketEvent, UserID string) {
	for client := range pool.Clients {
		if client.UserID == UserID {
			select {
			case client.Send <- payload:
			default:
				close(client.Send)
				delete(pool.Clients, client)
			}
		}
	}
}
