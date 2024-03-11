package service

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RoomID   string `json:"roomid"`
	Conn     *websocket.Conn
	Message  chan *Message
}
type Message struct {
	RoomID   string `json:"roomid"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()
	for {
		msg, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(msg)
	}
}

func (c *Client) Read(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error : %v", err)
			}
			break
		}
		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}
		log.Println(string(m))
		hub.Broadcast <- msg

	}

}
