import asyncio
import os
from dataclasses import dataclass
from pathlib import Path
from typing import Dict, Any, Optional, List
from uuid import UUID

import httpx
import socketio
from dotenv import load_dotenv


@dataclass
class AppConfig:
    auth_url: str = "http://localhost:5000/api/auth"
    room_url: str = "http://localhost:5000/api"
    chat_url: str = "http://localhost:5000/api"
    ws_api_url: str = "http://localhost:5000/api"
    ws_url: str = "http://localhost:5000"


class ApiClient:
    def __init__(self, base_url: str):
        self._client = httpx.AsyncClient(base_url=base_url, timeout=httpx.Timeout(30.0, connect=10.0))

    def set_auth(self, token: str) -> None:
        self._client.headers["Authorization"] = f"Bearer {token}"

    async def post(self, path: str, json: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        response = await self._client.post(path, json=json)
        response.raise_for_status()
        return response.json()

    async def post_multipart(self, path: str, data: Dict[str, Any]) -> Dict[str, Any]:
        response = await self._client.post(path, data=data)
        response.raise_for_status()
        return response.json()

    async def get(self, path: str) -> Dict[str, Any]:
        response = await self._client.get(path)
        response.raise_for_status()
        return response.json()

    async def patch(self, path: str, json: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        response = await self._client.patch(path, json=json)
        response.raise_for_status()
        return response.json()

    async def delete(self, path: str) -> Dict[str, Any]:
        response = await self._client.delete(path)
        response.raise_for_status()
        # DELETE may return 204 No Content (empty body)
        if response.status_code == 204:
            return {"success": True, "message": "Deleted successfully"}
        return response.json()

    async def close(self) -> None:
        await self._client.aclose()


class UserSession:
    def __init__(self, config: AppConfig, email: str, password: str):
        self.email = email
        self.password = password
        self.config = config
        self.auth_client = ApiClient(config.auth_url)
        self.room_client = ApiClient(config.room_url)
        self.chat_client = ApiClient(config.chat_url)
        self.ws_api_client = ApiClient(config.ws_api_url)
        self.sio = socketio.AsyncClient()
        
        # Event tracking for message operations
        self.new_message_received: asyncio.Event = asyncio.Event()
        self.message_updated_received: asyncio.Event = asyncio.Event()
        self.message_deleted_received: asyncio.Event = asyncio.Event()
        
        self.latest_message: Optional[Dict[str, Any]] = None
        self.latest_updated_message: Optional[Dict[str, Any]] = None
        self.deleted_message_id: Optional[str] = None

    async def setup(self) -> None:
        auth_data = await self.auth_client.post("/sign-in", {
            "email": self.email,
            "password": self.password
        })
        assert auth_data["success"] is True, f"Auth failed: {auth_data}"
        token = auth_data["data"]["accessToken"]

        self.auth_client.set_auth(token)
        self.room_client.set_auth(token)
        self.chat_client.set_auth(token)
        self.ws_api_client.set_auth(token)

    async def get_ws_token(self) -> str:
        res = await self.ws_api_client.get("/ws-token")
        assert res["success"] is True, f"WS token retrieval failed: {res}"
        return res["data"]

    async def connect_ws(self, chat_id: str) -> None:
        """Connect to WebSocket and setup event listeners"""
        ws_token = await self.get_ws_token()

        @self.sio.on("NEW_MESSAGE")
        async def on_new_message(data: Dict[str, Any]) -> None:
            print(f"[{self.email}] Received NEW_MESSAGE: {data}")
            self.latest_message = data
            self.new_message_received.set()

        @self.sio.on("MESSAGE_UPDATED")
        async def on_message_updated(data: Dict[str, Any]) -> None:
            print(f"[{self.email}] Received MESSAGE_UPDATED: {data}")
            self.latest_updated_message = data
            self.message_updated_received.set()

        @self.sio.on("MESSAGE_DELETED")
        async def on_message_deleted(data: Dict[str, Any]) -> None:
            print(f"[{self.email}] Received MESSAGE_DELETED: {data}")
            self.deleted_message_id = data.get("id")
            self.message_deleted_received.set()

        url = f"{self.config.ws_url}?token={ws_token}&chat_id={chat_id}"
        await self.sio.connect(url, transports=["websocket"])
        print(f"[{self.email}] Connected to WebSocket for chat {chat_id}")

    async def cleanup(self) -> None:
        try:
            if self.sio.connected:
                await self.sio.disconnect()
        except Exception:
            pass
        
        try:
            await self.auth_client.close()
        except Exception:
            pass
        
        try:
            await self.room_client.close()
        except Exception:
            pass
        
        try:
            await self.chat_client.close()
        except Exception:
            pass
        
        try:
            await self.ws_api_client.close()
        except Exception:
            pass


async def run_e2e_test() -> None:
    """
    E2E test for messages flow:
    1. Login with two users
    2. Create room
    3. Retrieve chat-id from room
    4. User 2 retrieves room to become a member
    5. Both retrieve ws-token
    6. Connect to Socket.IO with chat id
    7. Subscribe for NEW_MESSAGE, MESSAGE_UPDATED, MESSAGE_DELETED events
    8. Send message and verify both users receive it
    9. Retrieve list of messages with pagination
    10. Update message and verify everyone receives the event
    11. Delete message and verify everyone receives the event
    """
    # Load environment variables from .env.test.local
    env_file = Path(__file__).parent / ".env.test.local"
    load_dotenv(env_file)
    
    config = AppConfig()
    
    user1_email = os.getenv("USER1_EMAIL")
    user1_password = os.getenv("USER1_PASSWORD")
    user2_email = os.getenv("USER2_EMAIL")
    user2_password = os.getenv("USER2_PASSWORD")
    
    missing_vars = []
    if not user1_email:
        missing_vars.append("USER1_EMAIL")
    if not user1_password:
        missing_vars.append("USER1_PASSWORD")
    if not user2_email:
        missing_vars.append("USER2_EMAIL")
    if not user2_password:
        missing_vars.append("USER2_PASSWORD")
    
    if missing_vars:
        raise RuntimeError(f"Missing env vars: {', '.join(missing_vars)}")
    
    # Initialize user sessions
    user1 = UserSession(config, user1_email, user1_password)
    user2 = UserSession(config, user2_email, user2_password)
    
    created_room_id: Optional[str] = None
    created_chat_id: Optional[str] = None
    created_message_id: Optional[str] = None

    try:
        print("\n=== Step 1: Login with two users ===")
        await user1.setup()
        print(f"✓ User 1 ({user1_email}) logged in")
        
        await user2.setup()
        print(f"✓ User 2 ({user2_email}) logged in")

        print("\n=== Step 2: Create room ===")
        room_res = await user1.room_client.post_multipart("/rooms", {
            "name": "E2E Message Test Room",
            "categoryName": "Education",
            "isPrivate": "false"
        })
        assert room_res["success"] is True, f"Room creation failed: {room_res}"
        created_room_id = room_res["data"]["id"]
        print(f"✓ Room created: {created_room_id}")

        print("\n=== Step 3: Retrieve chat ID from room ===")
        room_info = await user1.room_client.get(f"/rooms/{created_room_id}")
        assert room_info["success"] is True, f"Room retrieval failed: {room_info}"
        created_chat_id = room_info["data"]["chatId"]
        print(f"✓ Chat ID retrieved: {created_chat_id}")

        print("\n=== Step 4: User 2 retrieves room to become a member ===")
        user2_room_res = await user2.room_client.get(f"/rooms/{created_room_id}")
        assert user2_room_res["success"] is True, f"Room retrieval by user2 failed: {user2_room_res}"
        print(f"✓ User 2 became a room member")

        print("\n=== Step 5 & 6: Retrieve WS tokens and connect to Socket.IO ===")
        await user1.connect_ws(created_chat_id)
        await user2.connect_ws(created_chat_id)
        print(f"✓ Both users connected to WebSocket")

        # Allow WebSocket connections to stabilize
        await asyncio.sleep(0.5)

        print("\n=== Step 7 & 8: Send message from user1 and verify user2 receives it ===")
        # Reset events
        user1.new_message_received.clear()
        user2.new_message_received.clear()

        send_res = await user1.chat_client.post(
            f"/chats/{created_chat_id}/messages",
            {"content": "Hello from User 1!"}
        )
        assert send_res["success"] is True, f"Message send failed: {send_res}"
        print(f"✓ Message post request successful")

        # Wait for both users to receive the NEW_MESSAGE event
        # Message ID comes from the WebSocket event, not from the HTTP response
        try:
            await asyncio.gather(
                asyncio.wait_for(user1.new_message_received.wait(), timeout=5),
                asyncio.wait_for(user2.new_message_received.wait(), timeout=5),
                return_exceptions=False
            )
            print(f"✓ Both users received NEW_MESSAGE event")
            assert user1.latest_message is not None, "User1 did not receive message event"
            assert user2.latest_message is not None, "User2 did not receive message event"
            assert user1.latest_message["content"] == "Hello from User 1!"
            assert user2.latest_message["content"] == "Hello from User 1!"
            
            # Extract message ID from WebSocket event payload
            created_message_id = user1.latest_message["id"]
            print(f"✓ Message created with ID (from WebSocket): {created_message_id}")
            print(f"✓ Message content verified for both users")
        except asyncio.TimeoutError:
            raise AssertionError("Timeout waiting for NEW_MESSAGE event")

        print("\n=== Step 9: Retrieve list of messages with pagination ===")
        messages_res = await user1.chat_client.get(f"/chats/{created_chat_id}/messages")
        assert messages_res["success"] is True, f"Message list retrieval failed: {messages_res}"
        messages_list = messages_res["data"]["items"]
        assert isinstance(messages_list, list), "Messages should be a list"
        assert len(messages_list) > 0, "Should have at least one message"
        
        # Verify the message we sent is in the list
        message_found = any(msg["id"] == created_message_id for msg in messages_list)
        assert message_found, f"Sent message {created_message_id} not found in messages list"
        print(f"✓ Message list retrieved: {len(messages_list)} message(s)")
        print(f"✓ Sent message found in list")

        print("\n=== Step 10: Update message and verify everyone receives the event ===")
        user1.message_updated_received.clear()
        user2.message_updated_received.clear()

        update_res = await user1.chat_client.patch(
            f"/chats/{created_chat_id}/messages/{created_message_id}",
            {"content": "Hello from User 1! (Updated)"}
        )
        assert update_res["success"] is True, f"Message update failed: {update_res}"
        print(f"✓ Message updated by User 1")

        try:
            await asyncio.gather(
                asyncio.wait_for(user1.message_updated_received.wait(), timeout=5),
                asyncio.wait_for(user2.message_updated_received.wait(), timeout=5),
                return_exceptions=False
            )
            print(f"✓ Both users received MESSAGE_UPDATED event")
            assert user1.latest_updated_message is not None
            assert user2.latest_updated_message is not None
            assert user1.latest_updated_message["content"] == "Hello from User 1! (Updated)"
            assert user2.latest_updated_message["content"] == "Hello from User 1! (Updated)"
            assert user1.latest_updated_message["isEdited"] is True
            assert user2.latest_updated_message["isEdited"] is True
            print(f"✓ Updated content and isEdited flag verified for both users")
        except asyncio.TimeoutError:
            raise AssertionError("Timeout waiting for MESSAGE_UPDATED event")

        print("\n=== Step 11: Delete message and verify everyone receives the event ===")
        user1.message_deleted_received.clear()
        user2.message_deleted_received.clear()

        delete_res = await user1.chat_client.delete(
            f"/chats/{created_chat_id}/messages/{created_message_id}"
        )
        # Delete returns 204 No Content, which may not parse as JSON
        print(f"✓ Message deleted by User 1")

        try:
            await asyncio.gather(
                asyncio.wait_for(user1.message_deleted_received.wait(), timeout=5),
                asyncio.wait_for(user2.message_deleted_received.wait(), timeout=5),
                return_exceptions=False
            )
            print(f"✓ Both users received MESSAGE_DELETED event")
            assert user1.deleted_message_id == created_message_id
            assert user2.deleted_message_id == created_message_id
            print(f"✓ Deleted message ID verified for both users")
        except asyncio.TimeoutError:
            raise AssertionError("Timeout waiting for MESSAGE_DELETED event")

        print("\n=== All tests passed! ===\n")

    finally:
        print("\n=== Cleanup ===")
        
        if created_room_id:
            try:
                await user1.room_client.delete(f"/rooms/{created_room_id}")
                print(f"✓ Room deleted: {created_room_id}")
            except httpx.HTTPStatusError:
                pass
        
        try:
            await user1.cleanup()
            print(f"✓ User 1 session cleaned up")
        except Exception:
            pass
        
        try:
            await user2.cleanup()
            print(f"✓ User 2 session cleaned up")
        except Exception:
            pass


if __name__ == "__main__":
    asyncio.run(run_e2e_test())
