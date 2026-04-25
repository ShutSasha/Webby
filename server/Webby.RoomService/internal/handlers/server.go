package handlers

import (
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"webby/internal/config"
	"webby/pkg/http/middleware/cors"
	loggerMw "webby/pkg/http/middleware/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewServer(config *config.Config, logger *slog.Logger, service Service) http.Handler {
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
	addRoutes(router, config, handler)

	return router
}
