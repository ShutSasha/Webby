# Webby.WsGateway

Standalone Socket.IO gateway that fronts all real-time traffic for the Webby platform.

## Architecture (one-liner)

```
                   gRPC (downstream)
client  <—— WS ——>  WsGateway  ─────────────►  ChatService / VotesService / RoomService
                       ▲
                       │ Redis Pub/Sub (upstream, channel "room:<uuid>")
                       │
              business services publish events
```

- **Downstream** — every action a client sends (`send_message`, `cast_vote`, …) is forwarded to the owning business service via gRPC. The gateway never touches the database.
- **Upstream** — business services publish JSON envelopes into Redis channels named `room:<room_id>`. The gateway re-broadcasts them to the matching Socket.IO room.
- **Activity worker** — every minute the gateway snapshots `{room_id -> [user_id]}` for currently connected sockets and ships it to RoomService via `AwardActiveUsers`. RoomService awards points and publishes `POINTS_AWARDED` back through Redis.

## Redis envelope

Every event published into `room:<room_id>` must look like:

```json
{
  "type": "MESSAGE_CREATED",     // or VOTE_UPDATED | POINTS_AWARDED
  "payload": { ... arbitrary JSON ... }
}
```

The gateway emits the inner `payload` to clients under an event name equal to `type`.

## Layout

```
cmd/server/main.go              wiring
config/config.yaml              env config
internal/config                 cleanenv loader
internal/clients                thin gRPC client wrappers (chat, votes, room)
internal/ws                     socket.io server + connection tracking
internal/redis                  Pub/Sub subscriber → ws broadcast
internal/worker                 activity ticker → RoomService.AwardActiveUsers
proto/{chat,votes,room}         .proto contracts
```

## Running

1. Generate proto code (output paths match `option go_package`):

   ```
   protoc --go_out=. --go-grpc_out=. \
          proto/chat/chat.proto proto/votes/votes.proto proto/room/room.proto
   ```

2. `go mod tidy && go run ./cmd/server`

The Socket.IO endpoint is served at `/socket.io/` on `cfg.http.port` (default `8090`).

## Connecting (client)

```
ws://localhost:8090/socket.io/?token=<JWT>&room_id=<UUID>
```

The gateway parses `sub` / `user_id` / `userId` / `uid` claims from the JWT (HMAC).

## Adding a new downstream event

1. Add a method on the relevant client in `internal/clients/`.
2. Register a handler in `internal/ws/server.go` (`registerHandlers`).
3. Make sure the corresponding business service publishes its result into `room:<id>` so the gateway can fan it out.
