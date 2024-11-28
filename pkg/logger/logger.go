package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var L *zap.Logger

func GetL() *zap.Logger {
	if L == nil {
		mainLogger := InitLogger()
		defer mainLogger.Sync()
		undo := zap.ReplaceGlobals(mainLogger)
		defer undo()
		L = zap.L()
	}
	return L
}

func LogError(msg string, err error) {
	l := GetL()
	l.Error(msg, zap.String("Err", err.Error()))
}

func ToZapLogLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func InitLogger() *zap.Logger {
	w := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		ToZapLogLevel("info"),
	)

	return zap.New(core)
}

func GetTestLogger() *zap.Logger {
	logConfig := zap.NewProductionConfig()
	logConfig.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	log, _ := logConfig.Build()
	return log
}
