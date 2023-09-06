package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (client *Client) ReadPump() {
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

		var socketEvent SocketEvent
		decoder := json.NewDecoder(bytes.NewReader(msg))
		decodeErr := decoder.Decode(&socketEvent)
		if decodeErr != nil {
			log.Println(decodeErr)
			return
		}

		HandleSocketPayloadEvents(client, socketEvent)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(((60 * time.Second) * 9) / 10)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case payload, ok := <-c.Send:
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(payload)
			finalPayload := reqBodyBytes.Bytes()

			fmt.Println("pushing message", string(finalPayload))

			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(finalPayload)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				json.NewEncoder(reqBodyBytes).Encode(<-c.Send)
				w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
