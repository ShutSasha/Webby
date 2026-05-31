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
    user_url: str = "http://localhost:5000/api/users"
    room_url: str = "http://localhost:5000/api"
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
    def __init__(self, config: AppConfig, email: str, password: str):
        self.email = email
        self.password = password
        self.config = config
        self.auth_client = ApiClient(config.auth_url)
        self.user_client = ApiClient(config.user_url)
        self.room_client = ApiClient(config.room_url)
        self.ws_api_client = ApiClient(config.ws_api_url)
        self.sio = socketio.AsyncClient()
        self.sync_completed = asyncio.Event()
        self.received_timecode = 0.0

    async def setup(self) -> None:
        auth_data = await self.auth_client.post("/sign-in", {
            "email": self.email,
            "password": self.password
        })
        token = auth_data["data"]["accessToken"]

        self.auth_client.set_auth(token)
        self.room_client.set_auth(token)
        self.ws_api_client.set_auth(token)

    async def get_ws_token(self) -> str:
        res = await self.ws_api_client.get("/ws-token")
        return res["data"]

    async def connect_ws(self, chat_id: str, room_id: str, report_value: int) -> None:
        ws_token = await self.get_ws_token()

        @self.sio.on("REPORT_TIMECODE")
        async def on_report_timecode(data: Dict[str, Any]) -> None:
            print(f"DEBUG REPORT_TIMECODE Payload: {data}")
            sync_id = data.get("synchronizeId", data.get("syncId"))
            await self.room_client.post(f"/rooms/{room_id}/sync/report", {
                "syncId": sync_id,
                "timecode": report_value
            })

        @self.sio.on("SYNCHRONIZE")
        async def on_synchronize(data: Dict[str, Any]) -> None:
            print(f"DEBUG SYNCHRONIZE Payload: {data}")
            self.received_timecode = data.get("timecode", 0)
            self.sync_completed.set()

        url = f"{self.config.ws_url}?token={ws_token}&chat_id={chat_id}"
        await self.sio.connect(url, transports=["websocket"])

    async def cleanup(self) -> None:
        if self.sio.connected:
            await self.sio.disconnect()
        
        await self.auth_client.close()
        await self.room_client.close()
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
    
    if not all([user1_email, user1_password, user2_email, user2_password]):
        raise ValueError("Missing required environment variables in .env.test.local file")
    
    user1 = UserSession(config, user1_email, user1_password)
    user2 = UserSession(config, user2_email, user2_password)

    try:
        await user1.setup()
        await user2.setup()

        room_res = await user1.room_client.post_multipart("/rooms", {
            "name": "E2E Test Room",
            "categoryName": "Education",
            "isPrivate": "false"
        })
        room_id = room_res["data"]["id"]

        room_info = await user1.room_client.get(f"/rooms/{room_id}")
        chat_id = room_info["data"]["chatId"]

        await user1.connect_ws(chat_id, room_id, 100)
        await user2.connect_ws(chat_id, room_id, 200)

        await user1.room_client.post(f"/rooms/{room_id}/sync")

        await asyncio.gather(
            asyncio.wait_for(user1.sync_completed.wait(), timeout=10),
            asyncio.wait_for(user2.sync_completed.wait(), timeout=10)
        )

        assert user1.received_timecode == 200
        assert user2.received_timecode == 200

        await user1.room_client.delete(f"/rooms/{room_id}")

    finally:
        await user1.cleanup()
        await user2.cleanup()

if __name__ == "__main__":
    asyncio.run(run_e2e_test())