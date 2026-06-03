---
name: create-e2e-test
description: Generates Python asyncio end-to-end (e2e) tests for backend APIs using httpx (and optionally socketio). Use this when asked to write or scaffold new e2e tests.
argument-hint: "[functionality to test] [optional edge cases]"
user-invocable: true
---

# Python E2E Test Generator

This skill helps you generate end-to-end tests for specific backend functionality, strictly following the project's established testing architecture.

## When to use this skill

Use this skill when you need to:
- Create new end-to-end tests for API endpoints.
- Write integration tests involving single or multiple user sessions.
- Scaffold WebSocket event testing scenarios (only if explicitly required by the flow).

## Reference Template

Before generating code, carefully review the [reference template](./synchronize.py) to understand the required structure, client wrappers, and session management. Adapt the template based on whether WebSockets are actually needed for the requested test.

## Architectural Requirements

Whenever you write an e2e test, you MUST strictly follow these rules:

1. **Environment & Auth Setup:** Always load variables from `.env.test.local` using `python-dotenv`. When authentication is required for the flow, you MUST extract test user credentials (e.g., emails, passwords) exclusively from these environment variables for however many users are needed. Fail early if required variables are missing.
2. **Client Wrappers:** - ALWAYS use `httpx.AsyncClient` for HTTP REST requests.
   - *OPTIONAL:* Only use `socketio.AsyncClient` if the specific test scenario involves WebSockets. Do not include it by default.
3. **Session Isolation:** Encapsulate user state within a `UserSession` class. Each actor (e.g., User1, User2) must have its own isolated instance containing its own clients and tokens.
4. **Lifecycle Hooks:** - Implement `async def setup()` to handle authentication and set Bearer tokens on the HTTP clients.
   - Implement `async def cleanup()` to close all HTTP clients and disconnect WebSockets (if they were initialized). 
5. **Strict Execution Flow & Infrastructure Cleanup (The Finally Block):** - The main `run_e2e_test()` function MUST use a `try...finally` block.
   - **Resource Teardown (Anti-Spam):** Inside the `try` block, track all infrastructure resources created during the test (e.g., `room_id`, `chat_id`). You MUST explicitly call the corresponding DELETE endpoints (or cleanup API methods) to completely remove these created resources from the system before the test ends.
   - **Connection Teardown:** You MUST guarantee that the `cleanup()` method for every initialized user session is called inside the `finally` block, ensuring network connections are freed regardless of assertion failures.
6. **Async Coordination:** If WebSockets are used, use `asyncio.Event()` for waiting on specific events.

## Process

1. Identify the target functionality and user flow described by the developer.
2. Determine how many isolated user sessions are required and whether WebSockets are necessary for this specific flow.
3. Generate a complete, self-contained Python script implementing the scenario.
4. Ensure all test logic and assertions (`assert`) are placed securely in the `try` block, and all cleanups are strictly in the `finally` block.