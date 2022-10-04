package logger

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("can't init logger", err)
	}

	Logger = logger.Sugar()
}
