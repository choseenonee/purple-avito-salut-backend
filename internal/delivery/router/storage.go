package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"template/internal/delivery/handlers"
	"template/internal/repository"
	"template/internal/service"
	"template/pkg/config"
	"template/pkg/log"
)

func RegisterStorageRouter(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer, urls []string) *gin.RouterGroup {
	storageRouter := r.Group("/storage")

	matrixRepo := repository.InitMatrixRepo(db, viper.GetInt(config.MaxOnPage))

	matrixService := service.InitUpdateService(matrixRepo)
	storageUpdateHandlers := handlers.InitMUpdateHandler(matrixService, tracer, urls)

	storageRouter.POST("/send", storageUpdateHandlers.PrepareAndSendStorage)
	storageRouter.POST("/switch", storageUpdateHandlers.SwitchStorageToNext)
	storageRouter.GET("/current", storageUpdateHandlers.GetCurrentStorage)

	return storageRouter
}
