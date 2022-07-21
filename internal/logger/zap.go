package logger

import (
	"errors"
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

var ErrInvalidLogLevel = errors.New("zap_logger: invalid log level")

func NewZap(level string) (*zap.Logger, error) {
	zapLogLevel, ok := logLevelsDict[strings.ToLower(level)]
	if !ok {
		return nil, fmt.Errorf("%w=s", ErrInvalidLogLevel)
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
		return nil, err
	}

	return logger, nil
}

func newZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "channel",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
