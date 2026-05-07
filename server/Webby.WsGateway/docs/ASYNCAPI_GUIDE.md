# Webby WebSocket Gateway - AsyncAPI Documentation

## Overview

The Webby WebSocket Gateway provides real-time chat messaging capabilities using Socket.IO. This document describes the WebSocket API, including connection requirements, events, and message formats.

## Getting Started

### Prerequisites

1. **Authentication Token**: Obtain a WebSocket token from the REST API endpoint `/api/ws-token`
2. **Chat ID**: The UUID of the chat room you want to join

### Connection URL

**Production:**
```
wss://api.webby.com:443/?token=<TOKEN>&chat_id=<CHAT_ID>
```

**Development:**
```
ws://localhost:8080/?token=<TOKEN>&chat_id=<CHAT_ID>
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

## API Events

### 1. Socket Connection (Implicit)

**Description:** Establishing a WebSocket connection to the server.

**Required Query Parameters:**
- `token` (string, JWT): WebSocket authentication token
- `chat_id` (string, UUID): The chat room ID to join

**Example Request:**
```
ws://localhost:8080/?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9&chat_id=550e8400-e29b-41d4-a716-446655440000
```

**Success Response:**
- Connection established
- Client automatically joins the chat room specified by `chat_id`
- Client can now send and receive messages

**Error Responses:**

| Error | Description |
|-------|-------------|
| `token and chat_id are required` | Missing token or chat_id parameter |
| `invalid chat_id` | chat_id is not a valid UUID |
| `unauthorized` | Token is invalid or expired |
| `internal server error` | Server encountered an error validating the token |

---

### 2. Send Message Event

**Event Name:** `send_message`

**Direction:** Client → Server

**Description:** Send a message to the current chat room. The message is persisted in the database and broadcasted to all clients connected to the same room.

**Request Payload:**
```json
{
  "content": "Your message here"
}
```

**Request Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `content` | string | Yes | Message content (1-5000 characters) |

**Example Request:**
```javascript
socket.emit('send_message', {
  content: 'Hello everyone!'
}, (response) => {
  console.log(response);
});
```

---

### 3. Send Message Response

**Event Name:** `send_message` (acknowledged via callback)

**Direction:** Server → Client

**Description:** Server response confirming whether the message was successfully saved.

**Response Payload (Success):**
```json
{
  "ok": true,
  "message_id": "660e8400-e29b-41d4-a716-446655440001"
}
```

**Response Payload (Failure):**
```json
{
  "ok": false,
  "error": "content is required"
}
```

**Response Fields:**
| Field | Type | Description |
|-------|------|-------------|
| `ok` | boolean | Operation status (true = success, false = failure) |
| `message_id` | string (UUID) | ID of the saved message (only if ok=true) |
| `error` | string | Error description (only if ok=false) |

**Possible Errors:**
| Error | Cause |
|-------|-------|
| `content is required` | Request payload is missing the content field |
| `failed to save message` | Server failed to persist the message to database |

**Example Response Handling:**
```javascript
socket.emit('send_message', 
  { content: 'Hello!' }, 
  (response) => {
    if (response.ok) {
      console.log('Message saved with ID:', response.message_id);
    } else {
      console.error('Error:', response.error);
    }
  }
);
```

---

### 4. Broadcast Message Event

**Event Name:** `message`

**Direction:** Server → Client (broadcast to all in room)

**Description:** Server broadcasts a message to all clients connected to the chat room after it has been successfully persisted.

**Broadcast Payload:**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "chat_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "770e8400-e29b-41d4-a716-446655440002",
  "content": "Hello everyone!",
  "created_at": "2026-05-05T10:30:00Z"
}
```

**Broadcast Fields:**
| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Unique message identifier |
| `chat_id` | string (UUID) | Chat room identifier |
| `user_id` | string (UUID) | Sender's user identifier |
| `content` | string | Message content |
| `created_at` | string (ISO 8601) | Server timestamp |

**Example Listener:**
```javascript
socket.on('message', (msg) => {
  console.log(`${msg.user_id}: ${msg.content}`);
  console.log(`Sent at: ${msg.created_at}`);
});
```

---

## Complete Usage Example

### JavaScript/TypeScript Client

```javascript
import { io } from 'socket.io-client';

// Step 1: Get WebSocket token from REST API
const response = await fetch('/api/ws-token', {
  headers: {
    'Authorization': 'Bearer <YOUR_JWT_TOKEN>'
  }
});
const { data: wsToken } = await response.json();

// Step 2: Connect to WebSocket
const socket = io('ws://localhost:8080', {
  query: {
    token: wsToken,
    chat_id: '550e8400-e29b-41d4-a716-446655440000'
  },
  transports: ['websocket']
});

// Step 3: Handle connection events
socket.on('connect', () => {
  console.log('Connected to WebSocket');
});

socket.on('connect_error', (error) => {
  console.error('Connection failed:', error);
});

// Step 4: Listen for incoming messages
socket.on('message', (msg) => {
  console.log(`${msg.user_id}: ${msg.content}`);
});

// Step 5: Send a message
socket.emit('send_message', 
  { content: 'Hello from the client!' },
  (response) => {
    if (response.ok) {
      console.log('Message sent successfully:', response.message_id);
    } else {
      console.error('Failed to send message:', response.error);
    }
  }
);

// Step 6: Handle disconnection
socket.on('disconnect', () => {
  console.log('Disconnected from WebSocket');
});
```

---

## Authentication & Security

### Token Generation
Tokens are generated via the REST API endpoint and are **one-time use**. Each WebSocket connection requires a fresh token.

### Token Validation
- Tokens are validated using JWT with the configured secret
- Invalid or expired tokens result in connection rejection
- User identity is verified during connection

### Room Isolation
- Each chat room is isolated - clients can only see messages from their connected room
- Users can connect to multiple rooms in separate WebSocket connections

---

## Error Handling

### Connection Errors
```javascript
socket.on('connect_error', (error) => {
  if (error.message === 'unauthorized') {
    // Token is invalid - request a new one
  } else if (error.message === 'invalid chat_id') {
    // Chat ID format is invalid
  } else {
    // Other connection errors
  }
});
```

### Send Message Errors
```javascript
socket.emit('send_message', 
  { content: 'Hello' },
  (response) => {
    if (!response.ok) {
      switch(response.error) {
        case 'content is required':
          console.error('Message content cannot be empty');
          break;
        case 'failed to save message':
          console.error('Database error - message not saved');
          break;
        default:
          console.error('Unknown error:', response.error);
      }
    }
  }
);
```

---

## Timeout & Limits

- **Call Timeout:** 5 seconds for message save operations
- **Content Length:** 1-5000 characters
- **Connection Timeout:** Standard Socket.IO defaults

---

## Related Resources

- **OpenAPI Spec:** See [oas.json](./oas.json) for REST API documentation
- **Socket.IO Docs:** https://socket.io/docs/
- **AsyncAPI Spec:** See [asyncapi.json](./asyncapi.json) for complete AsyncAPI specification

---

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-05-05 | Initial release with send_message event and socket connection |

