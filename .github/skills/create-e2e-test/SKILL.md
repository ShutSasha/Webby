---
name: create-e2e-test
description: Generates Python asyncio end-to-end (e2e) tests for backend APIs using httpx (and optionally socketio). Use this when asked to write or scaffold new e2e tests.
argument-hint: "[functionality to test] [optional edge cases]"
user-invocable: true
---

# Python E2E Test Generator

This skill helps you generate end-to-end tests for specific backend functionality, strictly following the project's established testing architecture.

## Reference Template
Before generating code, carefully review the [reference template](./synchronize.py) to understand the expected `ApiClient` and `UserSession` implementations.

## Implementation Checklist & Rules

When writing an e2e test, you MUST execute the generation by following these exact sequential steps:

### Step A: Environment & Initialization
1. Load variables from `.env.test.local` using `python-dotenv`.
2. **Strict Env Naming:** Use explicitly named variables like `USER1_EMAIL`, `USER1_PASSWORD`, `USER2_EMAIL`, etc., up to the number of users required by the flow.
3. **Fail-Fast:** If any required variable is missing, raise a RuntimeError listing them.
   *(Example: `raise RuntimeError(f"Missing env vars: {missing}")`)*.

### Step B: Session & Client Configuration
1. **Timeouts:** Instantiate `httpx.AsyncClient` with an explicit timeout to prevent hanging tests *(Example: `timeout=httpx.Timeout(30.0, connect=10.0)`)*.
2. **Socket.IO Rule:** ONLY include `socketio.AsyncClient` and WebSocket event logic IF the user explicitly mentions real-time events, socket messages, or WebSocket endpoints. Otherwise, omit them entirely.
3. **UserSession API:** Ensure the `UserSession` class contains distinct `setup()` and `cleanup()` methods, managing its own token injection for its HTTP clients.

### Step C: Test Execution (The `try` block)
1. The main logic MUST reside inside a `try` block.
2. **Resource Tracking:** When a POST request creates a resource (e.g., room, category), store its ID in a variable (e.g., `created_room_id = response["data"]["id"]`) so it can be accessed in the `finally` block later.
3. **API Response Envelope (Strict Parsing):** ALL backend endpoints return a standard envelope. You MUST NOT assume the raw entity is at the root level.
   - For general requests, the payload is ALWAYS inside the `"data"` key.
   - For paginated requests, the array is ALWAYS inside `"data"]["items"]`.
4. **Diagnostic Assertions:** When asserting success, include the full response in the error message to aid debugging.
   *(Example: `assert response["success"] is True, f"API error: {response}"`)*.

### Step D: Teardown & Cleanup (The `finally` block)
1. The `finally` block MUST handle both infrastructure cleanup (database entities) and connection teardown.
2. **Idempotent Resource Deletion:** For every resource ID tracked in Step C, execute a DELETE HTTP request. Wrap each DELETE call in its own `try...except httpx.HTTPStatusError: pass` block so a failure in one does not prevent others from deleting.
3. **Session Teardown:** Finally, call `.cleanup()` on every initialized `UserSession` object. Wrap each cleanup call in its own `try...except Exception: pass` block to guarantee all clients are closed.