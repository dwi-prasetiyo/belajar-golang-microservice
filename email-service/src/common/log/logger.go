package log

import "go.uber.org/zap"

var Logger *zap.Logger

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	
	Logger = logger
}
