package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const dateTimeLayout = "2006-01-02 15:04:05"

type Logger interface {
	Info(...interface{})
	Error(...interface{})
	Panic(...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
	Panicw(string, ...interface{})
	Sync() error
}

func New() Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(dateTimeLayout)
	zapLogger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			zap.NewAtomicLevel(),
		),
	)
	return zapLogger.Sugar()
}
