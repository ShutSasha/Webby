---
name: Go Clean Code Rules
description: Enforces Clean Code, SOLID, DRY, KISS, and self-documenting principles for backend development.
applyTo: '**/*.go' 
---

# Global Go Backend Engineering Standards

- Never generate `.md` files for documentation. 
*Reason:* Documentation must be delivered as plain text in the chat, and all code must live strictly in `.go` files to avoid repository clutter.

- Write straightforward, simple code and implement only what is strictly required (KISS & YAGNI). 
*Reason:* Over-engineering, speculative features, or "clever" tricks make the codebase harder to read and maintain for the rest of the team.

- Extract repetitive logic into reusable functions or shared packages (DRY). 
*Reason:* Duplicated code increases the risk of bugs when logic needs to be updated in multiple places.

- Define interfaces directly at the consumer's place of use. 
*Reason:* Go uses implicit interfaces. Defining the interface where it is used (rather than where it is implemented) prevents tight coupling and circular dependencies.
*Avoid:*
```go
// In package "store"
type UserRepository interface {
    GetUser(id string) (*User, error)
}
type dbRepo struct{} // Implementation bundled with interface

```

*Prefer:*
```go
// In package "service" (the consumer)
type UserFetcher interface {
    GetUser(id string) (*User, error)
}
func ProcessUser(f UserFetcher, id string) { ... }

```

- Do not use inline comments or docstrings to explain *what* code is doing.
*Reason:* Comments often drift out of sync with the code. The code itself should be highly expressive and self-documenting.

Extract complex logic into heavily descriptive variables or functions instead of writing comments.
*Reason:* Descriptive naming natively documents the code and makes unit testing easier.
*Avoid:*

```go
// Check if the user is old enough and active
if u.Age >= 18 && u.Status == "active" { ... }

```

*Prefer:*

```go
if isAdultAndActive(u) { ... }

```

- Never hardcode magic values (e.g., timeout durations, URLs, retry limits) directly within business logic.
*Reason:* Hardcoded values make the code rigid, difficult to test, and prone to breaking across different environments.

*Avoid:*
```go
client := http.Client{Timeout: 30 * time.Second}

```

*Prefer:*
```go
const defaultClientTimeout = 30 * time.Second
client := http.Client{Timeout: defaultClientTimeout}

```

- Break down large, monolithic structs and packages into smaller components with a Single Responsibility.
*Reason:* "God objects" create high coupling and low cohesion, making the system fragile and difficult to test in isolation.
