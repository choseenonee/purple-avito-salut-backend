package main

import (
	"fmt"
	"github.com/spf13/viper"
	"template/internal/delivery"
	"template/internal/delivery/middleware"
	"template/pkg/config"
	"template/pkg/database"
	"template/pkg/log"
	"template/pkg/tracing"
)

const serviceName = "admin-panel"

func main() {
	logger, loggerInfoFile, loggerErrorFile := log.InitLogger()
	defer loggerInfoFile.Close()
	defer loggerErrorFile.Close()

	logger.Info("Logger Initialized")

	config.InitConfig()
	logger.Info("Config Initialized")

	jaegerURL := fmt.Sprintf("http://%v:%v/api/traces", viper.GetString(config.JaegerHost), viper.GetString(config.JaegerPort))
	tracer := tracing.InitTracer(jaegerURL, serviceName)
	logger.Info("Tracer Initialized")

	db := database.GetDB()
	logger.Info("Database Initialized")

	mdw := middleware.InitMiddleware(logger)

	urls := []string{viper.GetString(config.ChildPortOne), viper.GetString(config.ChildPortTwo)}

	delivery.Start(
		db,
		logger,
		tracer,
		mdw,
		urls,
	)
}
