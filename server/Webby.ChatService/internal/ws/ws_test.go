package ws_test

import (
	"testing"

	"webby-chat/internal/ws"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHubManager_GetOrCreateHub(t *testing.T) {
	manager := ws.NewHubManager()
	chatId := uuid.New()

	hub1 := manager.GetOrCreateHub(chatId)
	require.NotNil(t, hub1)
	require.Equal(t, chatId, hub1.ChatID())

	hub2 := manager.GetOrCreateHub(chatId)
	require.Equal(t, hub1, hub2)

	chatId2 := uuid.New()
	hub3 := manager.GetOrCreateHub(chatId2)
	require.NotEqual(t, hub1, hub3)
}

func TestHubManager_RemoveHub(t *testing.T) {
	manager := ws.NewHubManager()
	chatId := uuid.New()

	_ = manager.GetOrCreateHub(chatId)
	require.NotNil(t, manager.GetHub(chatId))

	manager.RemoveHub(chatId)
	require.Nil(t, manager.GetHub(chatId))
}

func TestHubManager_ClientCount(t *testing.T) {
	manager := ws.NewHubManager()
	require.Equal(t, 0, manager.ClientCount())
}

func TestHub_RegisterUnregister(t *testing.T) {
	chatId := uuid.New()
	hub := ws.NewHub(chatId)
	go hub.Run()
	defer hub.Stop()

	userId := uuid.New()
	client := ws.NewClient(userId, chatId)

	hub.Register(client)

	require.Eventually(t, func() bool {
		return hub.ClientCount() == 1
	}, 1e9, 1e6)

	hub.Unregister(client)

	require.Eventually(t, func() bool {
		return hub.ClientCount() == 0
	}, 1e9, 1e6)
}

func TestHub_Broadcast(t *testing.T) {
	chatId := uuid.New()
	hub := ws.NewHub(chatId)
	go hub.Run()
	defer hub.Stop()

	userId := uuid.New()
	client := ws.NewClient(userId, chatId)
	hub.Register(client)

	require.Eventually(t, func() bool {
		return hub.ClientCount() == 1
	}, 1e9, 1e6)

	hub.Broadcast("test:event", map[string]string{"hello": "world"})

	msg := <-client.Send()
	require.Equal(t, "test:event", msg.Event)
	require.Equal(t, map[string]string{"hello": "world"}, msg.Data)
}

func TestClient_Creation(t *testing.T) {
	userId := uuid.New()
	chatId := uuid.New()
	client := ws.NewClient(userId, chatId)

	require.Equal(t, userId, client.UserID)
	require.Equal(t, chatId, client.ChatID)
	require.NotNil(t, client.Send())
}
