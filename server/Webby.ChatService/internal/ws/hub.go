package ws

import (
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	chatId     uuid.UUID
	clients    map[uuid.UUID]*Client
	mu         sync.RWMutex
	broadcast  chan *BroadcastMessage
	register   chan *Client
	unregister chan *Client
	stop       chan struct{}
}

type BroadcastMessage struct {
	Event string
	Data  any
}

func NewHub(chatId uuid.UUID) *Hub {
	return &Hub{
		chatId:     chatId,
		clients:    make(map[uuid.UUID]*Client),
		broadcast:  make(chan *BroadcastMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		stop:       make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, client.UserID)
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- msg:
				default:
				}
			}
			h.mu.RUnlock()

		case <-h.stop:
			return
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(event string, data any) {
	h.broadcast <- &BroadcastMessage{
		Event: event,
		Data:  data,
	}
}

func (h *Hub) Stop() {
	close(h.stop)
}

func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

func (h *Hub) ChatID() uuid.UUID {
	return h.chatId
}
