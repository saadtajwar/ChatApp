package websocket

import (
	"fmt"
	"log"
)

func (client *Client) Read() {
	defer func() {
		client.Pool.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		payload := Payload{
			UserID:   client.UserID,
			Username: client.Username,
			Message:  string(msg),
		}
		message := SocketEvent{EventName: "Broadcast", EventPayload: payload}
		client.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
