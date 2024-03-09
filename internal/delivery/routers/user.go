package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"template/internal/delivery/handlers"
	"template/internal/repository/user"
	"template/internal/service"
	"template/pkg/auth"
	"template/pkg/database/cached"
	"template/pkg/log"
)

func RegisterUserRouter(userRouter *gin.RouterGroup, db *sqlx.DB, session cached.Session, jwt auth.JWTUtil, logger *log.Logs, tracer trace.Tracer) {
	userRepo := user.InitUserRepo(db)

	userService := service.InitUserService(userRepo, session, jwt, logger)
	userHandler := handlers.InitUserHandler(userService, session, tracer)

	userRouter.GET("/me", userHandler.GetMe)
	userRouter.GET("/delete", userHandler.Delete)
}
