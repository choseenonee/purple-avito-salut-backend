package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"template/internal/delivery/docs"
	"template/internal/delivery/middleware"
	"template/internal/delivery/routers"
	"template/pkg/auth"
	"template/pkg/database/cached"
	"template/pkg/log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start(db *sqlx.DB, logger *log.Logs, session cached.Session, tracer trace.Tracer) {
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	jwtUtils := auth.InitJWTUtil()
	middlewareStruct := middleware.InitMiddleware(logger, jwtUtils, session)

	r.Use(middlewareStruct.CORSMiddleware())

	routers.InitRouting(r, db, logger, middlewareStruct, jwtUtils, session, tracer)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(fmt.Sprintf("error running client: %v", err.Error()))
	}
}
