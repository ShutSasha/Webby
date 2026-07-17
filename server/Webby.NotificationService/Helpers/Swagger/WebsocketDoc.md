# Notification Service Real-Time API

Real-time notification service for the platform. Provides SignalR/WebSocket events for receiving push notifications and unread message counts in real-time.

## Connection Details

Connect to the notifications hub using the following endpoint:

`ws://localhost:5000/hubs/notifications`

**Authentication:** Authentication is required and utilizes a secure one-time ticket system.

1. **Obtain a Ticket:** Make an authenticated request to:
   `GET /api/notification/one-time-ticket`
   *(Extract the ticket string from the API response).*

2. **Establish Connection:** Connect to the WebSocket hub by passing the retrieved ticket via the `ticket` query parameter:
   `ws://localhost:5000/hubs/notifications?ticket=<YOUR_ONE_TIME_TICKET>`

---

## Server → Client Events

The client should subscribe to the following events to receive updates from the server.

| Event | Payload Type   | Description                                                                                                                                                          |
| :--- |:---------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `UpdateUnreadNotificationsCount` | `Integer`      | Receives the up-to-date count of unread notifications for the connected user. Usually emitted upon connection or after a new notification is created.                |
| `AuthError` | `String`       | Receives an error message if the connection is rejected (e.g., "Unauthorized: Token is missing or expired"). The connection is aborted immediately after this event. |
| `ReceiveNotification` | `Notification` | Broadcasts a new push notification to the user.                                                                                                                      |
| `AccountSuspended` | `Guid`         | Broadcasts a system notification of user block.                                                                                                                      |

---

## Data Models

### Notification Object

The payload structure received during the `ReceiveNotification` event.

```json
{
  "notificationId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "userId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "title": "New Subscriber!",
  "message": "Someone just followed your channel.",
  "targetType": "System" | "Video" | "Playlist" | "Room" | "User", 
  "targetIdentifier": "some-action-target-id",
  "createdAt": "2023-10-27T10:00:00Z",
  "notificationStatus": "Read" | "Unread"
}