# Achievement Service

The Achievement Service is responsible for tracking user progress and unlocking milestones based on cross-service events. It operates asynchronously by consuming telemetry data from a shared Redis Stream.

## 🐳 Infrastructure Initialization (Redis)

To start the Redis message broker locally without Docker Compose, ensure your Docker daemon is running and execute the following command in your terminal. This will spin up a lightweight Redis container in detached mode and expose the default port (no password required).

```bash
docker run -d \
  --name platform-redis \
  -p 6379:6379 \
  redis:7.2-alpine \
  redis-server --appendonly yes
```

## Event-Driven Architecture (Redis Streams)

Other microservices must publish their domain events to the Redis broker to register user progress. The Achievement Service listens to these events and increments user progress accordingly.

* **Stream Name:** `events:platform`

### Base Payload Structure
Every event pushed to the stream MUST adhere to the following base JSON structure:

```json
{
  "eventType": "string (e.g., video_watched, interaction_like)",
  "Payload": {
     "userId": "Interaction user id (uuid format)",
     "isIncrementOperation": "Specify type of storing progress value in redis (incrementing value or reset to absolute value)",
     "value": "integer interactions value"
  },
  "Timestamp": "ISO 8601 UTC DateTime"
}
```

