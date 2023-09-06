package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

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
		fmt.Println("register statement")
		registerEvent := SocketEvent{
			EventName: socketEvent.EventName,
			EventPayload: Payload{
				UserID:   client.UserID,
				Username: client.Username,
				Users:    GetAllConnectedUsers(client.Pool),
			},
		}
		BroadcastMessageToAll(client.Pool, registerEvent)
	case "disconnect":
		fmt.Println("disconnect statement")
		disconnectEvent := SocketEvent{
			EventName: socketEvent.EventName,
			EventPayload: Payload{
				UserID:   client.UserID,
				Username: client.Username,
				Users:    GetAllConnectedUsers(client.Pool),
			},
		}
		BroadcastMessageToAll(client.Pool, disconnectEvent)
	case "message":
		fmt.Println("direct message")
		socketEventResponse := SocketEvent{
			EventName: "message response",
			EventPayload: Payload{
				Username: socketEvent.EventPayload.Username,
				UserID:   socketEvent.EventPayload.UserID,
				Message:  socketEvent.EventPayload.Message,
			},
		}
		EmitToSpecificClient(client.Pool, socketEventResponse, socketEventResponse.EventPayload.UserID)
	}

}

func EmitToSpecificClient(pool *Pool, payload SocketEvent, UserID string) {
	wg := new(sync.WaitGroup)

	for client := range pool.Clients {
		if client.UserID == UserID {
			c := client
			wg.Add(1)
			go func(wg *sync.WaitGroup, c *Client) {
				defer wg.Done()
				SendMessageToChannel(c, payload)
			}(wg, c)
		}
	}

	fmt.Println("USERID of the client getting the message is ", UserID)
	fmt.Println("EmitToSpecificClient :: Message Broadcasted, waiting for WaitGroup")
	wg.Wait()

}

func BroadcastMessageToAll(pool *Pool, socketEvent SocketEvent) {

	wg := new(sync.WaitGroup)

	for client := range pool.Clients {

		c := client
		wg.Add(1)
		go func(wg *sync.WaitGroup, c *Client) {
			defer wg.Done()
			SendMessageToChannel(c, socketEvent)
		}(wg, c)
		// select {
		// case client.Send <- socketEvent:
		// 	// default:
		// 	// 	close(client.Send)
		// 	// 	delete(pool.Clients, client)
		// }
	}
	fmt.Println("BroadcastMessageToAll :: Message Broadcasted, waiting for WaitGroup")
	wg.Wait()

}

func SendMessageToChannel(client *Client, socketEvent SocketEvent) {
	fmt.Println("The userid of the client for the SendMessageToChannel ", client.UserID)
	client.Send <- socketEvent
}

func GetAllConnectedUsers(pool *Pool) []User {
	var users []User
	for client := range pool.Clients {
		user := User{Username: client.Username, UserID: client.UserID}
		users = append(users, user)
	}
	return users
}
