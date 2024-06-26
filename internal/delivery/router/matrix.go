package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	handlers "template/internal/delivery/handlers"
	"template/internal/repository"
	"template/internal/service"
	"template/pkg/config"
	"template/pkg/log"
)

func RegisterMatrixUser(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) *gin.RouterGroup {
	matrixRouter := r.Group("/matrix")

	matrixRepo := repository.InitMatrixRepo(db, viper.GetInt(config.MaxOnPage))

	matrixService := service.InitMatrixService(matrixRepo)
	matrixHandlers := handlers.InitMatrixHandler(matrixService, tracer)

	matrixRouter.POST("/create", matrixHandlers.CreateMatrix)
	matrixRouter.POST("/create_no_parent", matrixHandlers.CreateMatrixWithoutParent)
	matrixRouter.GET("/get_difference", matrixHandlers.GetDifference)
	matrixRouter.GET("/get_matrices_by_duration", matrixHandlers.GetMatricesByDuration)
	matrixRouter.GET("/get_matrix", matrixHandlers.GetMatrix)
	matrixRouter.PUT("/get_history", matrixHandlers.GetHistory)
	matrixRouter.PUT("/get_tendency", matrixHandlers.GetTendency)
	matrixRouter.GET("/get_matrix_pages", matrixHandlers.GetMatrixPages)

	return matrixRouter
}
