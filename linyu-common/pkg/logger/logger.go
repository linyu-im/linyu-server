package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLog() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}
