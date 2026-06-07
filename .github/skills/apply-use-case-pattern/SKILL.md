---
name: apply-use-case-pattern
description: Refactors large God-object services into isolated, single-responsibility Use Cases (Interactors) following Clean Architecture principles. Use this when asked to split a large service, fix a file with too many dependencies, or apply the use case pattern.
argument-hint: target file or service struct
---

# Use Case Pattern Refactoring Guide

This skill helps you refactor bloated, multi-responsibility services ("God Objects") into isolated, highly cohesive Use Cases (also known as Interactors or Commands). 

## When to use this skill
- When a service struct has five or more dependencies (>= 5) in its constructor.
- Also consider refactoring when unit tests require mocking unrelated interfaces or when method cohesion is low (e.g., methods touch unrelated aggregates). Do not rely solely on the '>=5' rule.
- When a user asks to apply "Clean Architecture", "CQRS", or "Use Case" patterns to a specific domain.
- When separating core domain logic (e.g., CRUD operations) from specialized operations (e.g., WebSocket sync, real-time events).

## Code Generation Constraints (STRICT)
- **No Comments (Self-Documenting Code):** Do not write explanatory comments in the generated code. The code must be clean, readable, and self-explanatory through proper variable, struct, and method naming (following clean code principles).
- **No Stubs or Placeholders:** Never generate placeholders, `TODO`s, or pseudo-code (e.g., `// log error here`, `// implement logic here`). You must write the actual, fully functional implementation. If logging or error handling is required, write the actual logger calls and return statements.
- **Do Not Hardcode:** Extract configurable values into variables, constants, or configuration parameters.

## Refactoring Process

Precedence: (1) Numbered steps, (2) User explicit preferences, (3) Best Practices, (4) Default YAGNI constraint (unless user requested generalization in explicit phrasing).

1. **Analyze Dependencies & Boundaries:**
   *Checkpoint A: Static Grouping*
   - 1.1 List methods and dependencies.
   - 1.2 Build a dependency-method matrix.
   - 1.3 Compute overlap = `|intersection(deps(methodA), deps(methodB))| / |union(deps(methodA), deps(methodB))|`. Group methods when overlap >= 0.60.
   
   *Checkpoint B: Behavioral Checks*
   - 1.4 For overlapping methods, check for shared transactions or ordered event pipelines. If the original method relied on a DB transaction spanning multiple responsibilities, preserve the transaction boundary: either keep those operations in a single Use Case, or inject and propagate an explicit transaction/session object. Document where the transaction starts and ends. If the platform/framework does not support injecting/propagating transactions, keep the original monolith method and mark it as 'transactional' in migration notes; provide tests proving behavior remains atomic.
   - 1.5 Write unit tests asserting: (a) each extracted Use Case's Execute covers the original method behavior for sample inputs; (b) side effects (DB writes/events) occur as before; (c) transaction boundaries are preserved. Run and pass these tests before extraction. If Checkpoint B fails, do not extract; document the reason.

2. **Extract the Use Case:** Create a new struct named after the specific business action (e.g., `RoomCreator`, `PlaybackSynchronizer`, `UserAuthenticator`).
3. **Isolate Dependencies:** Inject *only* the interfaces required for this specific use case into its constructor. Do not pass the entire repository layer if only one method is needed.
   - *Cross-cutting concerns:* For cross-cutting concerns (logging, metrics, auth, tracing) inject minimal abstractions (e.g., `Logger`, `MetricsRecorder`, `AuthContext`) into the Use Case only if the Use Case directly needs them. Otherwise keep these in the delivery layer and pass context objects; document each decision.
   - *Globals/Cyclic Imports:* If package-level globals or cyclic imports prevent injecting minimal interfaces, provide concrete steps: (a) introduce an adapter interface around the global, (b) extract a thin interface facade in the same package, or (c) refactor package boundaries. If you encounter cyclic imports or cannot create a minimal interface due to package coupling, stop extraction and return a report with (1) files involved, (2) suggested package splits, and (3) a temporary adapter pattern to maintain functionality until refactor is applied.
4. **Implement Execute:** Use `Execute` as the canonical public entrypoint. Put helper logic in unexported (private) methods. Only expose additional public methods if they represent distinct use cases. Constructors should return concrete pointers (e.g., `*RoomCreator`) unless the caller requires mocking — in that case, provide a constructor that returns an interface with only the `Execute` method. Document visibility decisions in a migration note.
5. **Update Delivery Layer (Handlers):** Generate the updated HTTP/gRPC handler code demonstrating how to inject and call the new Use Cases. Apply the Interface Segregation Principle (ISP) by defining minimal interfaces for each Use Case directly inside the handler package. Do not inject the concrete Use Case struct into the handler.

## Best Practices (KISS, YAGNI, DRY, SOLID)
- **Accept Interfaces, Return Structs (Default):** Declare minimal interfaces in the same package as the use case, in a file named `<usecase>_interfaces.go` (or inline if very short). Do not place them in the repository package. Constructors return structs by default, returning interfaces only when explicitly needed for mocking (as defined in Step 4).
- **Fail-Fast:** Validate inputs and contexts at the very beginning of the `Execute` method.
- **Do not over-engineer (YAGNI):** Treat the user as requesting a generalization if they use phrases like "make reusable", "generalize", "add abstraction", or ask for a library/module-level API. Otherwise apply YAGNI. Do not create a complex generic Use Case interface or bus dispatcher unless explicitly requested. Direct method calls are preferred for simplicity (KISS).

## Example Pattern (Go)

**Before (Bloated Service):**
```go
type RoomService struct {
    roomRepo       RoomRepository
    fileRepo       FileRepository
    categoryClient CategoryClient
    publisher      EventPublisher
    timecodesRepo  TimecodesRepo
}

func (s *RoomService) Create(ctx context.Context, room *Room) error { return nil }
func (s *RoomService) ReportTimecode(ctx context.Context, code int) error { return nil }
```
**After (Refactored Use Cases & Delivery Layer):**
File: internal/services/room_creator.go
```go
package services

import "context"

type RoomCreatorRepo interface {
    Create(ctx context.Context, room *Room) error
}
type FileUploader interface {
    Upload(ctx context.Context, file []byte) error
}

type RoomCreator struct {
    repo     RoomCreatorRepo
    uploader FileUploader
}

func NewRoomCreator(repo RoomCreatorRepo, uploader FileUploader) *RoomCreator {
    return &RoomCreator{repo: repo, uploader: uploader}
}

func (uc *RoomCreator) Execute(ctx context.Context, room *Room, file []byte) error {
    if err := uc.uploader.Upload(ctx, file); err != nil {
        return err
    }
    
    if err := uc.repo.Create(ctx, room); err != nil {
        return err
    }
    
    return nil
}
```

File: internal/handlers/handler.go
```go
...

func (h *handler) Create(c *gin.Context) {
    ...
    err := h.roomCreator.Execute(c.Request.Context(), &req.Room, req.File)

    ...
}
```
