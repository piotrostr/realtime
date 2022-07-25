package logger

import (
	"log"

	"go.uber.org/zap"
)

func Get() *zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	return sugar
}
