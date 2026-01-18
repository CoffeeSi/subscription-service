package logging

import (
	"go.uber.org/zap"
)

func NewLogging() (*zap.Logger, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return logger, nil
}
