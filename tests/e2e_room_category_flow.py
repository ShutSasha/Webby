import asyncio
import os
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Dict, Optional

import httpx
from dotenv import load_dotenv


@dataclass
class AppConfig:
    auth_url: str = "http://localhost:5000/api/auth"
    category_url: str = "http://localhost:5000/api"


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

    async def put(self, path: str, json: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        response = await self._client.put(path, json=json)
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
        self.category_client = ApiClient(config.category_url)

    async def setup(self) -> None:
        auth_data = await self.auth_client.post("/sign-in", {
            "email": self.email,
            "password": self.password,
        })

        token = auth_data["data"]["accessToken"]
        self.category_client.set_auth(token)

    async def cleanup(self) -> None:
        await self.auth_client.close()
        await self.category_client.close()


async def run_e2e_test() -> None:
    env_file = Path(__file__).parent / ".env.test.local"
    load_dotenv(env_file)

    user1_email = os.getenv("USER1_EMAIL")
    user1_password = os.getenv("USER1_PASSWORD")

    if not all([user1_email, user1_password]):
        raise ValueError("Missing required environment variables in .env.test.local file")

    config = AppConfig()
    user1 = UserSession(config, user1_email, user1_password)

    created_category_name: Optional[str] = None
    updated_category_name = "Updated Category"

    try:
        await user1.setup()

        category_name = "Test Category"
        created_category_name = category_name

        create_response = await user1.category_client.post("/categories", {"name": category_name})
        assert create_response["success"] is True

        list_response = await user1.category_client.get("/categories")
        items = list_response["data"]["items"]
        assert category_name in items

        update_response = await user1.category_client.put(f"/categories/{category_name}", {"name": updated_category_name})
        assert update_response["success"] is True
        created_category_name = updated_category_name

        list_response = await user1.category_client.get("/categories")
        items = list_response["data"]["items"]
        assert updated_category_name in items
        assert category_name not in items

        delete_response = await user1.category_client.delete(f"/categories/{updated_category_name}")
        assert delete_response["success"] is True
        created_category_name = None

        list_response = await user1.category_client.get("/categories")
        items = list_response["data"]["items"]
        assert updated_category_name not in items

    finally:
        if created_category_name:
            try:
                await user1.category_client.delete(f"/categories/{created_category_name}")
            except httpx.HTTPStatusError:
                pass
        await user1.cleanup()


if __name__ == "__main__":
    asyncio.run(run_e2e_test())
