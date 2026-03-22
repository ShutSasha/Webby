package repository

import "errors"

var (
	ErrNotFound  = errors.New("resource not found")
	ErrConflict  = errors.New("resource already exists")
	ErrInvalidID = errors.New("invalid id")
)

// Custom errors for user-friendly messages
type RepositoryError struct {
	Code    string // Machine-readable code
	Message string // User-friendly message
	Err     error  // Original error
}

func (e *RepositoryError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// Room-specific errors
func ErrRoomNotFound(id string) *RepositoryError {
	return &RepositoryError{
		Code:    "ROOM_NOT_FOUND",
		Message: "Room not found: " + id,
		Err:     ErrNotFound,
	}
}

func ErrRoomCreationFailed(reason string) *RepositoryError {
	return &RepositoryError{
		Code:    "ROOM_CREATION_FAILED",
		Message: "Failed to create room: " + reason,
	}
}

func ErrRoomUpdateFailed(reason string) *RepositoryError {
	return &RepositoryError{
		Code:    "ROOM_UPDATE_FAILED",
		Message: "Failed to update room: " + reason,
	}
}

func ErrRoomDeletionFailed(reason string) *RepositoryError {
	return &RepositoryError{
		Code:    "ROOM_DELETION_FAILED",
		Message: "Failed to delete room: " + reason,
	}
}

func ErrRoomsFetchFailed(reason string) *RepositoryError {
	return &RepositoryError{
		Code:    "ROOMS_FETCH_FAILED",
		Message: "Failed to fetch rooms: " + reason,
	}
}

// Category-specific errors
func ErrCategoryNotFound(id string) *RepositoryError {
	return &RepositoryError{
		Code:    "CATEGORY_NOT_FOUND",
		Message: "Category not found: " + id,
		Err:     ErrNotFound,
	}
}

func ErrCategoryCreationFailed(reason string) *RepositoryError {
	return &RepositoryError{
		Code:    "CATEGORY_CREATION_FAILED",
		Message: "Failed to create category: " + reason,
	}
}

func ErrCategoryUpdateFailed(reason string) *RepositoryError {
	return &RepositoryError{
		Code:    "CATEGORY_UPDATE_FAILED",
		Message: "Failed to update category: " + reason,
	}
}

func ErrCategoryDeletionFailed(reason string) *RepositoryError {
	return &RepositoryError{
		Code:    "CATEGORY_DELETION_FAILED",
		Message: "Failed to delete category: " + reason,
	}
}

func ErrCategoriesFetchFailed(reason string) *RepositoryError {
	return &RepositoryError{
		Code:    "CATEGORIES_FETCH_FAILED",
		Message: "Failed to fetch categories: " + reason,
	}
}

// Generic database errors
func ErrDatabaseConnection(reason string) *RepositoryError {
	return &RepositoryError{
		Code:    "DB_CONNECTION_ERROR",
		Message: "Database connection error: " + reason,
	}
}
