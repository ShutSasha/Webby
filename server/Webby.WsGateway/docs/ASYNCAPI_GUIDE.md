# Webby WebSocket Gateway - AsyncAPI Documentation

## Overview

The Webby WebSocket Gateway provides real-time chat messaging capabilities using Socket.IO. This document describes the WebSocket API, including connection requirements, events, and message formats.

## Getting Started

### Prerequisites

1. **Authentication Token**: Obtain a WebSocket token from the REST API endpoint `/api/ws-token`
2. **Chat ID**: The UUID of the chat room you want to join

### Connection URL

**Development:**
```
ws://localhost:5000/?token=<TOKEN>&chat_id=<CHAT_ID>
```

### Connection Flow

```
1. Client requests WebSocket token from /api/ws-token (REST API)
   ↓
2. Server returns one-time token
   ↓
3. Client initiates WebSocket connection with token and chat_id query parameters
   ↓
4. Server validates token and chat_id
   ↓
5. On success: Client joins the chat room and can send/receive messages
   On failure: Connection is rejected with error message
```

## API Events to subscribe to

### 1. QUEUE_UPDATED

**Description:** Notifies all room members that the room queue has been updated.

**Payload**
```json
{
  "position": 1
}
```

**Possible use-case:** After receiving the position of an added or deleted item, the client can check if that position is currently in view. If so, the client should refresh the queue list using `GET /api/rooms/{id}/queue`


