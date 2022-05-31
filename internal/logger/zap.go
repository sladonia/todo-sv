package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevelsDict = map[string]zapcore.Level{
	"debug":   zapcore.DebugLevel,
	"info":    zapcore.InfoLevel,
	"warning": zapcore.WarnLevel,
	"error":   zapcore.ErrorLevel,
	"panic":   zapcore.PanicLevel,
	"fatal":   zapcore.FatalLevel,
}

func MustCreateZapLogger(level string) *zap.Logger {
	zapLogLevel, ok := logLevelsDict[strings.ToLower(level)]
	if !ok {
		panic(fmt.Sprintf("unsupported log level=%s", level))
	}

	encConfig := newZapEncoderConfig()

	zapConfig := zap.Config{
		Level: zap.NewAtomicLevelAt(zapLogLevel),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func newZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
