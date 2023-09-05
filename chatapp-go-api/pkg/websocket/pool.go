package websocket

import "fmt"

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			HandleUserRegister(pool, client)
		case client := <-pool.Unregister:
			HandleUserDisconnect(pool, client)
		}
	}
}

func HandleUserRegister(pool *Pool, client *Client) {
	pool.Clients[client] = true
	fmt.Println("Size of Connection Pool: ", len(pool.Clients))
	joinEvent := SocketEvent{
		EventName: "register",
		EventPayload: Payload{
			UserID: client.UserID,
		},
	}
	HandleSocketPayloadEvents(client, joinEvent)
}

func HandleUserDisconnect(pool *Pool, client *Client) {
	_, ok := pool.Clients[client]
	if ok {
		delete(pool.Clients, client)
		fmt.Println("Size of Connection Pool: ", len(pool.Clients))
		close(client.Send)
		leaveEvent := SocketEvent{
			EventName: "disconnect",
			EventPayload: Payload{
				UserID: client.UserID,
			},
		}
		HandleSocketPayloadEvents(client, leaveEvent)
	}
}
