package list

import (
	"context"
	"log/slog"
	"net/http"
	_ "webby/internal/handlers/docs"
	errorWrapper "webby/internal/handlers/errors"
	"webby/internal/handlers/responses"
	"webby/internal/services"
	"webby/pkg/http/render"

	"github.com/google/uuid"
)

type QueueLister interface {
	GetQueue(ctx context.Context, roomId, userId uuid.UUID) ([]services.QueueItemEnriched, error)
}

func New(logger *slog.Logger, lister QueueLister) http.Handler {
	return errorWrapper.MakeHandler(logger, listQueue(logger, lister))
}

// listQueue godoc
// @Summary      Get room queue
// @Description  Retrieve the playback queue for a room with enriched media data.
// @Description  Videos include title and thumbnail. Playlists include title, thumbnail, and child videos.
// @Description  **Path Parameter Validation:**
// @Description  * `id`: must be a valid UUID v4 format (room ID)
// @Description  **Access Rules:**
// @Description  * Only **room members** can view the queue
// @Description  **Security Note:**
// @Description  * This action **requires authentication**
// @Tags         Room Queue
// @Produce      json
// @Param        id path string true "Room ID (UUID v4 format)"
// @Success      200 {object} docs.ApiResponse[[]docs.QueueItemDetailResponse] "Queue successfully retrieved"
// @Failure      400 {object} docs.Error400Response "Invalid UUID format in path"
// @Failure      401 {object} docs.Error401Response "Missing or invalid authentication token"
// @Failure      403 {object} docs.Error403Response "Access denied - only room members can view the queue"
// @Failure      500 {object} docs.Error500Response "Internal server error"
// @Security     BearerAuth
// @Router       /rooms/{id}/queue [get]
func listQueue(logger *slog.Logger, lister QueueLister) errorWrapper.APIFunc {
	log := logger.With(slog.String("operation", "httpserver.roomqueues.list"))

	type videoChild struct {
		Id         uuid.UUID `json:"id"`
		Title      string    `json:"title"`
		Thumbnail  string    `json:"thumbnail"`
		VideoUrl   string    `json:"videoUrl"`
		PreviewUrl string    `json:"previewUrl"`
	}

	type queueItemResponse struct {
		Id         uuid.UUID    `json:"id"`
		EntityId   uuid.UUID    `json:"entityId"`
		EntityType string       `json:"entityType"`
		Title      string       `json:"title"`
		Thumbnail  string       `json:"thumbnail"`
		VideoUrl   string       `json:"videoUrl"`
		PreviewUrl string       `json:"previewUrl"`
		IsActive   bool         `json:"isActive"`
		IsFolder   bool         `json:"isFolder"`
		Children   []videoChild `json:"children,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		roomIdStr := r.PathValue("id")
		roomId, err := uuid.Parse(roomIdStr)
		if err != nil {
			log.Error("Invalid room id", slog.String("err", err.Error()))
			return responses.NewValidationError("Validation failed", map[string]string{
				"id": "the room id format is not valid",
			})
		}

		userIdStr := r.Context().Value("userID").(string)
		userId, _ := uuid.Parse(userIdStr)

		items, err := lister.GetQueue(r.Context(), roomId, userId)
		if err != nil {
			log.Error("List queue error", slog.String("err", err.Error()))
			return responses.NewApiError("List queue error", err)
		}

		result := make([]queueItemResponse, 0, len(items))
		for _, item := range items {
			entry := queueItemResponse{
				Id:         item.Id,
				EntityId:   item.EntityId,
				EntityType: item.EntityType,
				Title:      item.Title,
				Thumbnail:  item.Thumbnail,
				VideoUrl:   item.VideoUrl,
				PreviewUrl: item.PreviewUrl,
				IsActive:   item.IsActive,
				IsFolder:   item.IsFolder,
			}
			if len(item.Children) > 0 {
				children := make([]videoChild, 0, len(item.Children))
				for _, c := range item.Children {
					children = append(children, videoChild{
						Id:         c.Id,
						Title:      c.Title,
						Thumbnail:  c.Thumbnail,
						VideoUrl:   c.VideoUrl,
						PreviewUrl: c.PreviewUrl,
					})
				}
				entry.Children = children
			}
			result = append(result, entry)
		}

		render.Encode(w, r, http.StatusOK, responses.ApiResponse[[]queueItemResponse]{
			Success: true,
			Message: "Queue retrieved",
			Data:    &result,
		})
		return nil
	}
}
