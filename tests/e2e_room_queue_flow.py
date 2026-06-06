import asyncio
import os
from dataclasses import dataclass
from pathlib import Path
from typing import Dict, Any, Optional

import httpx
import socketio
from dotenv import load_dotenv


@dataclass
class AppConfig:
    auth_url: str = "http://localhost:5000/api/auth"
    room_url: str = "http://localhost:5000/api"
    video_url: str = "http://localhost:5000/api"
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
        self.video_client = ApiClient(config.video_url)
        self.ws_api_client = ApiClient(config.ws_api_url)
        self.sio = socketio.AsyncClient()
        self.queue_updated = asyncio.Event()

    async def setup(self) -> None:
        auth_data = await self.auth_client.post("/sign-in", {
            "email": self.email,
            "password": self.password
        })
        token = auth_data["data"]["accessToken"]

        self.auth_client.set_auth(token)
        self.room_client.set_auth(token)
        self.video_client.set_auth(token)
        self.ws_api_client.set_auth(token)

    async def get_ws_token(self) -> str:
        res = await self.ws_api_client.get("/ws-token")
        return res["data"]

    async def connect_ws(self, chat_id: str) -> None:
        ws_token = await self.get_ws_token()

        @self.sio.on("QUEUE_UPDATED")
        async def on_queue_updated(data: Dict[str, Any]) -> None:
            print(f"DEBUG QUEUE_UPDATED Payload: {data}")
            self.queue_updated.set()

        url = f"{self.config.ws_url}?token={ws_token}&chat_id={chat_id}"
        await self.sio.connect(url, transports=["websocket"])

    async def cleanup(self) -> None:
        if self.sio.connected:
            await self.sio.disconnect()
        
        await self.auth_client.close()
        await self.room_client.close()
        await self.video_client.close()
        await self.ws_api_client.close()


async def run_e2e_test() -> None:
    env_file = Path(__file__).parent / ".env.test.local"
    load_dotenv(env_file)
    
    config = AppConfig()
    
    user1_email = os.getenv("USER1_EMAIL")
    user1_password = os.getenv("USER1_PASSWORD")
    
    if not all([user1_email, user1_password]):
        raise ValueError("Missing required environment variables (USER1_EMAIL, USER1_PASSWORD) in .env.test.local file")
    
    user1 = UserSession(config, user1_email, user1_password)

    try:
        await user1.setup()

        room_res = await user1.room_client.post("/rooms", {
            "name": "E2E Room Queue Test",
            "categoryName": "Education",
            "isPrivate": False
        })
        room_id = room_res["data"]["id"]
        print(f"Created room: {room_id}")

        room_info = await user1.room_client.get(f"/rooms/{room_id}")
        chat_id = room_info["data"]["chatId"]
        print(f"Retrieved chatId: {chat_id}")

        await user1.connect_ws(chat_id)
        print("Connected to WebSocket")

        videos_res = await user1.video_client.get("/videos/search")
        videos = videos_res.get("data", {}).get("items", [])
        if not videos:
            raise ValueError("No videos available in VideoService")
        video_id_1 = videos[0]["videoId"]
        video_id_2 = videos[1]["videoId"] if len(videos) > 1 else video_id_1
        print(f"Retrieved videoIds: {video_id_1}, {video_id_2}")

        user1.queue_updated.clear()
        add_res_1 = await user1.room_client.post(f"/rooms/{room_id}/queue", {
            "videoId": video_id_1
        })
        print(f"Added first video to queue: {add_res_1}")
        
        await asyncio.wait_for(user1.queue_updated.wait(), timeout=10)
        print("Received QUEUE_UPDATED event after adding first video")
        
        queue_res_1 = await user1.room_client.get(f"/rooms/{room_id}/queue")
        queue_1 = queue_res_1.get("data", {}).get("items", [])
        print(f"Queue after first video: {queue_1}")
        assert len(queue_1) == 1, "Queue should contain 1 item after adding first video"

        user1.queue_updated.clear()
        add_res_2 = await user1.room_client.post(f"/rooms/{room_id}/queue", {
            "videoId": video_id_2
        })
        print(f"Added second video to queue: {add_res_2}")
        
        await asyncio.wait_for(user1.queue_updated.wait(), timeout=10)
        print("Received QUEUE_UPDATED event after adding second video")
        
        queue_res_2 = await user1.room_client.get(f"/rooms/{room_id}/queue")
        queue_2 = queue_res_2.get("data", {}).get("items", [])
        print(f"Queue after second video: {queue_2}")
        assert len(queue_2) == 2, "Queue should contain 2 items after adding second video"

        item_id = queue_2[1]["id"]
        print(f"Second video itemId: {item_id}")

        user1.queue_updated.clear()
        activate_res = await user1.room_client.patch(f"/rooms/{room_id}/queue/{item_id}/activate", {})
        print(f"Activated second video: {activate_res}")
        
        await asyncio.wait_for(user1.queue_updated.wait(), timeout=10)
        print("Received QUEUE_UPDATED event after activating second video")
        
        queue_res_3 = await user1.room_client.get(f"/rooms/{room_id}/queue")
        queue_3 = queue_res_3.get("data", {}).get("items", [])
        print(f"Queue after activation: {queue_3}")
        
        activated_item = next((item for item in queue_3 if item["id"] == item_id), None)
        assert activated_item is not None, "Activated item should still be in queue"
        assert activated_item.get("isActive", False) is True, "Second video should have isActive=true"
        print("✓ Second video is now active")

        user1.queue_updated.clear()
        delete_res = await user1.room_client.delete(f"/rooms/{room_id}/queue/{item_id}")
        print(f"Deleted second video from queue: {delete_res}")
        
        await asyncio.wait_for(user1.queue_updated.wait(), timeout=10)
        print("Received QUEUE_UPDATED event after deleting second video")
        
        queue_res_4 = await user1.room_client.get(f"/rooms/{room_id}/queue")
        queue_4 = queue_res_4.get("data", {}).get("items", [])
        print(f"Queue after deletion: {queue_4}")
        
        deleted_item = next((item for item in queue_4 if item["id"] == item_id), None)
        assert deleted_item is None, "Deleted item should no longer be in queue"
        print("✓ Deleted item successfully removed from queue")

        print("\n✅ All assertions passed!")

    finally:
        try:
            await user1.room_client.delete(f"/rooms/{room_id}")
            print(f"Cleaned up room: {room_id}")
        except Exception as e:
            print(f"Failed to cleanup room: {e}")
        
        await user1.cleanup()


if __name__ == "__main__":
    asyncio.run(run_e2e_test())
