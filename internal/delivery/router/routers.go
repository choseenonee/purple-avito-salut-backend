package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"template/pkg/log"
)

func InitRouting(r *gin.Engine, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer, urls []string) {
	_ = RegisterMatrixUser(r, db, logger, tracer)
	_ = RegisterStorageRouter(r, db, logger, tracer, urls)
}
