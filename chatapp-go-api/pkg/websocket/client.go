package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Username string
	Conn     *websocket.Conn
	Pool     *Pool
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (client *Client) Read() {
	defer func() {
		client.Pool.Unregister <- client
		client.Conn.Close()
	}()

	for {
		messageType, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(msg)}
		client.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
