package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	handlers "template/internal/delivery/handlers"
	"template/internal/repository"
	"template/internal/service"
	"template/pkg/log"
)

func RegisterMatrixUser(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	matrixRouter := r.Group("/matrix")

	matrixRepo := repository.InitMatrixRepo(db)

	matrixService := service.InitMatrixService(matrixRepo)
	matrixHandlers := handlers.InitMatrixHandler(matrixService, tracer)

	matrixRouter.GET("/get_difference", matrixHandlers.GetDifference)
	matrixRouter.POST("/create", matrixHandlers.CreateMatrix)
	matrixRouter.PUT("/get_history", matrixHandlers.GetHistory)
	matrixRouter.PUT("/get_tendency", matrixHandlers.GetTendency)

	return matrixRouter
}
