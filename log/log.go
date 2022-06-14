package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger() *zap.Logger {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// Default output log in stdout
	core := zapcore.NewCore(encoder, os.Stdout, zap.DebugLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	return logger
}
