package handlers

import (
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"webby-chat/internal/config"
	"webby-chat/pkg/http/middleware/cors"
	loggerMw "webby-chat/pkg/http/middleware/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	socketio "github.com/googollee/go-socket.io"
)

func NewServer(
	cfg *config.Config,
	logger *slog.Logger,
	service Service,
	socketServer *socketio.Server,
) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		v.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
			return strings.TrimSpace(fl.Field().String()) != ""
		})
	}

	router.Use(loggerMw.Logger(logger))
	router.Use(gin.Recovery())
	router.Use(cors.CORS())

	handler := New(service)
	addRoutes(router, cfg, handler)

	return router
}
