package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	handlers "template/internal/delivery/handelrs"
	"template/internal/repository"
	"template/internal/service"
	"template/pkg/log"
)

func RegisterMatrixUser(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	matrixRouter := r.Group("/matrix")

	matrixRepo := repository.InitMatrixRepo(db)

	matrixService := service.InitMatrixService(matrixRepo)
	matrixHandlers := handlers.InitMatrixHandler(matrixService, tracer)

	matrixRouter.POST("/create", matrixHandlers.CreateMatrix)
	matrixRouter.PUT("/get_history", matrixHandlers.GetHistory)

	return matrixRouter
}
