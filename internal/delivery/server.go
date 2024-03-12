package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"template/internal/delivery/docs"
	"template/internal/delivery/handlers"
	"template/internal/delivery/middleware"
	"template/internal/repository"
	service2 "template/internal/service"
	"template/pkg/log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start(db *sqlx.DB, logger *log.Logs, tracer trace.Tracer, middleware middleware.Middleware) {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(middleware.CORSMiddleware())

	repo := repository.InitRepository(db)

	// FIXME: matrix name чета сделать надо...
	service := service2.InitService(repo, "baseline_1710203720")

	handler := handlers.InitHandler(service, tracer)

	r.PUT("/price", handler.GetPrice)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
