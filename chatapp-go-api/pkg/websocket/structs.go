package websocket

import "github.com/gorilla/websocket"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan SocketEvent
}

type Client struct {
	Username string
	Conn     *websocket.Conn
	Pool     *Pool
	Send     chan SocketEvent
	UserID   string
}

type SocketEvent struct {
	EventName    string  `json:"eventname"`
	EventPayload Payload `json:"eventpayload"`
}

type Payload struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
