package main

import (
	"fmt"
	"github.com/spf13/viper"
	"template/internal/delivery"
	"template/pkg/config"
	"template/pkg/database"
	"template/pkg/database/cached"
	"template/pkg/log"
	"template/pkg/trace"
)

const serviceName = "gin"

func main() {
	logger, loggerInfoFile, loggerErrorFile := log.InitLogger()
	defer loggerInfoFile.Close()
	defer loggerErrorFile.Close()

	logger.Info("Logger Initialized")

	config.InitConfig()
	logger.Info("Config Initialized")

	jaegerURL := fmt.Sprintf("http://%v:%v/api/traces", viper.GetString(config.JaegerHost), viper.GetString(config.JaegerPort))
	tracer := trace.InitTracer(jaegerURL, serviceName)
	logger.Info("Tracer Initialized")

	db := database.GetDB()
	logger.Info("Database Initialized")

	redisSession := cached.InitRedis(tracer)
	logger.Info("Redis Initialized")

	delivery.Start(
		db,
		logger,
		redisSession,
		tracer,
	)

}
