import asyncio
import os
from dataclasses import dataclass, field
from pathlib import Path
from typing import Dict, Any, Optional, List

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

    async def delete(self, path: str) -> Dict[str, Any]:
        response = await self._client.delete(path)
        response.raise_for_status()
        return response.json()

    async def close(self) -> None:
        await self._client.aclose()


class UserSession:
    def __init__(self, config: AppConfig, email: str, password: str, user_label: str):
        self.email = email
        self.password = password
        self.user_label = user_label
        self.user_id = None  # Will be set during setup
        self.config = config
        self.auth_client = ApiClient(config.auth_url)
        self.room_client = ApiClient(config.room_url)
        self.chat_client = ApiClient(config.chat_url)
        self.ws_api_client = ApiClient(config.ws_api_url)
        self.sio = socketio.AsyncClient()
        self.received_messages: List[Dict[str, Any]] = []
        self.message_received = asyncio.Event()

    async def setup(self) -> None:
        auth_data = await self.auth_client.post("/sign-in", {
            "email": self.email,
            "password": self.password
        })
        token = auth_data["data"]["accessToken"]
        self.user_id = auth_data["data"]["user"]["userId"]

        self.auth_client.set_auth(token)
        self.room_client.set_auth(token)
        self.chat_client.set_auth(token)
        self.ws_api_client.set_auth(token)

    async def get_ws_token(self) -> str:
        res = await self.ws_api_client.get("/ws-token")
        return res["data"]

    async def connect_ws(self, chat_id: str) -> None:
        ws_token = await self.get_ws_token()

        @self.sio.on("NEW_MESSAGE")
        async def on_new_message(data: Dict[str, Any]) -> None:
            print(f"[{self.user_label}] Received NEW_MESSAGE event: {data}")
            self.received_messages.append(data)
            self.message_received.set()

        url = f"{self.config.ws_url}?token={ws_token}&chat_id={chat_id}"
        await self.sio.connect(url, transports=["websocket"])

    async def send_message(self, chat_id: str, content: str) -> Dict[str, Any]:
        response = await self.chat_client.post(f"/messages/{chat_id}", {
            "content": content
        })
        return response

    async def cleanup(self) -> None:
        if self.sio.connected:
            await self.sio.disconnect()
        
        await self.auth_client.close()
        await self.room_client.close()
        await self.chat_client.close()
        await self.ws_api_client.close()


async def run_e2e_test() -> None:
    # Load environment variables from .env.test.local
    env_file = Path(__file__).parent / ".env.test.local"
    load_dotenv(env_file)
    
    config = AppConfig()
    
    user1_email = os.getenv("USER1_EMAIL")
    user1_password = os.getenv("USER1_PASSWORD")
    user2_email = os.getenv("USER2_EMAIL")
    user2_password = os.getenv("USER2_PASSWORD")
    
    missing = []
    if not user1_email:
        missing.append("USER1_EMAIL")
    if not user1_password:
        missing.append("USER1_PASSWORD")
    if not user2_email:
        missing.append("USER2_EMAIL")
    if not user2_password:
        missing.append("USER2_PASSWORD")
    
    if missing:
        raise RuntimeError(f"Missing env vars: {missing}")
    
    user1 = UserSession(config, user1_email, user1_password, "user1_label")
    user2 = UserSession(config, user2_email, user2_password, "user2_label")
    
    created_room_id = None
    created_chat_id = None

    try:
        # Step 1: Login both users
        print("\n[Test] Step 1: Logging in both users...")
        await user1.setup()
        await user2.setup()
        print("[Test] Both users logged in successfully")

        # Step 2: Create a room and get its chat ID
        print("\n[Test] Step 2: Creating room and retrieving chat ID...")
        room_res = await user1.room_client.post("/rooms", {
            "name": "Message Flow Test Room",
            "categoryName": "Education",
            "isPrivate": False
        })
        assert room_res["success"] is True, f"Failed to create room: {room_res}"
        created_room_id = room_res["data"]["id"]
        print(f"[Test] Room created: {created_room_id}")

        # Get chat ID from room
        room_info = await user1.room_client.get(f"/rooms/{created_room_id}")
        assert room_info["success"] is True, f"Failed to get room info: {room_info}"
        created_chat_id = room_info["data"]["chatId"]
        print(f"[Test] Chat ID retrieved: {created_chat_id}")

        # Step 3: Add user2 as a room member by retrieving the room
        print("\n[Test] Step 3: Adding user2 as room member...")
        user2_room_info = await user2.room_client.get(f"/rooms/{created_room_id}")
        assert user2_room_info["success"] is True, f"Failed for user2 to get room info: {user2_room_info}"
        print("[Test] User2 is now a room member")

        # Step 4: Connect both users to the same chat via Socket.IO
        print("\n[Test] Step 4: Connecting both users to chat via Socket.IO...")
        await user1.connect_ws(created_chat_id)
        await user2.connect_ws(created_chat_id)
        print("[Test] Both users connected to WebSocket")

        # Small delay to ensure connection is stable
        await asyncio.sleep(0.5)

        # Step 5: Send message from user1 and verify both users receive it
        print("\n[Test] Step 5: User1 sends message...")
        user1_message = "Hello from User 1!"
        
        # Reset event flags BEFORE sending to ensure we catch the new message
        user1.message_received.clear()
        user2.message_received.clear()
        
        message1_res = await user1.send_message(created_chat_id, user1_message)
        assert message1_res["success"] is True, f"Failed to send message: {message1_res}"
        print(f"[Test] Message sent from user1: {user1_message}")

        # Wait for both users to receive the message
        print("[Test] Waiting for both users to receive NEW_MESSAGE event...")
        try:
            await asyncio.gather(
                asyncio.wait_for(user1.message_received.wait(), timeout=5),
                asyncio.wait_for(user2.message_received.wait(), timeout=5)
            )
        except asyncio.TimeoutError:
            raise AssertionError("Timeout waiting for NEW_MESSAGE event from one or both users")

        # Verify received messages
        assert len(user1.received_messages) > 0, "User1 did not receive any messages"
        assert len(user2.received_messages) > 0, "User2 did not receive any messages"

        received_msg_user1 = user1.received_messages[-1]
        received_msg_user2 = user2.received_messages[-1]

        print(f"[Test] User1 received message: {received_msg_user1}")
        print(f"[Test] User2 received message: {received_msg_user2}")

        # Verify message content
        assert received_msg_user1["content"] == user1_message, \
            f"Message content mismatch for user1: expected '{user1_message}', got '{received_msg_user1['content']}'"
        assert received_msg_user2["content"] == user1_message, \
            f"Message content mismatch for user2: expected '{user1_message}', got '{received_msg_user2['content']}'"
        assert received_msg_user1["sender"]["id"] == user1.user_id, \
            f"Sender ID mismatch: expected '{user1.user_id}', got '{received_msg_user1['sender']['id']}'"
        assert received_msg_user2["sender"]["id"] == user1.user_id, \
            f"Sender ID mismatch: expected '{user1.user_id}', got '{received_msg_user2['sender']['id']}'"

        print("[Test] ✓ Both users received message1 correctly")

        # Step 6: Send message from user2 and verify both users receive it
        print("\n[Test] Step 6: User2 sends message...")
        user2_message = "Hello from User 2!"
        
        # Reset event flags BEFORE sending to ensure we catch the new message
        user1.message_received.clear()
        user2.message_received.clear()
        
        message2_res = await user2.send_message(created_chat_id, user2_message)
        assert message2_res["success"] is True, f"Failed to send message: {message2_res}"
        print(f"[Test] Message sent from user2: {user2_message}")

        print("[Test] Waiting for both users to receive NEW_MESSAGE event...")
        try:
            await asyncio.gather(
                asyncio.wait_for(user1.message_received.wait(), timeout=5),
                asyncio.wait_for(user2.message_received.wait(), timeout=5)
            )
        except asyncio.TimeoutError:
            raise AssertionError("Timeout waiting for NEW_MESSAGE event from one or both users")

        received_msg_user1 = user1.received_messages[-1]
        received_msg_user2 = user2.received_messages[-1]

        print(f"[Test] User1 received message: {received_msg_user1}")
        print(f"[Test] User2 received message: {received_msg_user2}")

        # Verify message content
        assert received_msg_user1["content"] == user2_message, \
            f"Message content mismatch for user1: expected '{user2_message}', got '{received_msg_user1['content']}'"
        assert received_msg_user2["content"] == user2_message, \
            f"Message content mismatch for user2: expected '{user2_message}', got '{received_msg_user2['content']}'"
        assert received_msg_user1["sender"]["id"] == user2.user_id, \
            f"Sender ID mismatch: expected '{user2.user_id}', got '{received_msg_user1['sender']['id']}'"
        assert received_msg_user2["sender"]["id"] == user2.user_id, \
            f"Sender ID mismatch: expected '{user2.user_id}', got '{received_msg_user2['sender']['id']}'"

        print("[Test] ✓ Both users received message2 correctly")
        print("\n[Test] ✓✓✓ All tests passed! ✓✓✓")

    finally:
        print("\n[Test] Cleanup: Closing connections and deleting resources...")
        
        # Delete the created room
        if created_room_id:
            try:
                await user1.room_client.delete(f"/rooms/{created_room_id}")
                print(f"[Test] Room deleted: {created_room_id}")
            except httpx.HTTPStatusError:
                pass
        
        # Cleanup sessions
        try:
            await user1.cleanup()
        except Exception:
            pass
        
        try:
            await user2.cleanup()
        except Exception:
            pass
        
        print("[Test] Cleanup completed")


if __name__ == "__main__":
    asyncio.run(run_e2e_test())
