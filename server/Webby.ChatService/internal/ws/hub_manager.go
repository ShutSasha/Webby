package ws

import (
	"sync"

	"github.com/google/uuid"
)

type HubManager struct {
	hubs map[uuid.UUID]*Hub
	mu   sync.RWMutex
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[uuid.UUID]*Hub),
	}
}

func (m *HubManager) GetOrCreateHub(chatId uuid.UUID) *Hub {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hub, ok := m.hubs[chatId]; ok {
		return hub
	}

	hub := NewHub(chatId)
	m.hubs[chatId] = hub
	go hub.Run()
	return hub
}

func (m *HubManager) GetHub(chatId uuid.UUID) *Hub {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hubs[chatId]
}

func (m *HubManager) RemoveHub(chatId uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if hub, ok := m.hubs[chatId]; ok {
		hub.Stop()
		delete(m.hubs, chatId)
	}
}

func (m *HubManager) ClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, hub := range m.hubs {
		count += hub.ClientCount()
	}
	return count
}
