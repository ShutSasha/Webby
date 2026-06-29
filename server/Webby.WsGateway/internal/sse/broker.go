package sse

import (
	"sync"

	"github.com/google/uuid"
)

type broker struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]map[chan []byte]struct{}
}

func New() *broker {
	return &broker{
		clients: make(map[uuid.UUID]map[chan []byte]struct{}),
	}
}

func (b *broker) Subscribe(userID uuid.UUID) chan []byte {
	ch := make(chan []byte, 100)

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.clients[userID] == nil {
		b.clients[userID] = make(map[chan []byte]struct{})
	}
	b.clients[userID][ch] = struct{}{}

	return ch
}

func (b *broker) Unsubscribe(userID uuid.UUID, ch chan []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if clients, ok := b.clients[userID]; ok {
		delete(clients, ch)
		if len(clients) == 0 {
			delete(b.clients, userID)
		}
	}
	close(ch)
}

func (b *broker) Broadcast(userID uuid.UUID, event []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if clients, ok := b.clients[userID]; ok {
		for ch := range clients {
			select {
			case ch <- event:
			default:
			}
		}
	}
}
