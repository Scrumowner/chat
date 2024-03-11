package service

import (
	"fmt"
)

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"-"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			_, ok := h.Rooms[client.RoomID]
			if ok {
				r := h.Rooms[client.RoomID]
				if _, ok := r.Clients[client.ID]; !ok {
					r.Clients[client.ID] = client
				}

			}
		case client := <-h.Unregister:
			_, ok := h.Rooms[client.RoomID]
			if ok {
				_, ok := h.Rooms[client.RoomID].Clients[client.ID]
				if ok {
					if len(h.Rooms[client.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							RoomID:   client.RoomID,
							Username: client.Username,
							Content:  fmt.Sprint("%s leaft chat", client.Username),
						}
					}
					delete(h.Rooms[client.RoomID].Clients, client.ID)
					close(client.Message)
				}

			}
		case msg := <-h.Broadcast:
			_, ok := h.Rooms[msg.RoomID]
			if ok {
				for _, client := range h.Rooms[msg.RoomID].Clients {
					if msg.Username == client.Username {
						continue
					}
					client.Message <- msg
				}
			}
		}

	}
}

///
