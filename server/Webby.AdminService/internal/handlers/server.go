package handlers

import (
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"webby/admin-service/internal/config"
	"webby/admin-service/pkg/http/middleware/cors"
	loggerMw "webby/admin-service/pkg/http/middleware/logger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewServer(
	config *config.Config,
	logger *slog.Logger,
	complaintsService complaintsService,
	statsService statsService,
	healthCheckers ...HealthChecker,
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

	hc := NewHealthCheck(healthCheckers...)

	router.GET("/livez", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})
	router.GET("/readyz", func(c *gin.Context) {
		if err := hc.Check(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	handler := New(complaintsService, statsService)
	addRoutes(router, config, handler)

	return router
}
