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
	eventType := socketEvent.EventName
	switch eventType {
	case "register":
		fmt.Printf("register statement")
		registerEvent := SocketEvent{
			EventName: socketEvent.EventName,
			EventPayload: Payload{
				UserID: client.UserID,
				Users:  GetAllConnectedUsers(client.Pool),
			},
		}
		BroadcastMessageToAll(client.Pool, registerEvent)
	case "disconnect":
		fmt.Printf("disconnect statement")
		disconnectEvent := SocketEvent{
			EventName: socketEvent.EventName,
			EventPayload: Payload{
				UserID: client.UserID,
				Users:  GetAllConnectedUsers(client.Pool),
			},
		}
		BroadcastMessageToAll(client.Pool, disconnectEvent)
	case "message":
		fmt.Printf("direct message")
		socketEventResponse := SocketEvent{
			EventName:    "message response",
			EventPayload: socketEvent.EventPayload,
		}
		EmitToSpecificClient(client.Pool, socketEventResponse, socketEventResponse.EventPayload.UserID)
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

func BroadcastMessageToAll(pool *Pool, socketEvent SocketEvent) {
	for client, _ := range pool.Clients {
		select {
		case client.Send <- socketEvent:
		default:
			close(client.Send)
			delete(pool.Clients, client)
		}
	}
}

func GetAllConnectedUsers(pool *Pool) []User {
	var users []User
	for client := range pool.Clients {
		user := User{Username: client.Username, UserID: client.UserID}
		users = append(users, user)
	}
	return users
}
