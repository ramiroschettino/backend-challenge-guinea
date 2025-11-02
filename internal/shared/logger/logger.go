package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}

type zapLogger struct {
	logger *zap.Logger
}

func NewLogger(level, format string) (Logger, error) {
	var config zap.Config

	if format == "json" {
		config = zap.NewProductionConfig() 
	} else {
		config = zap.NewDevelopmentConfig()
	}

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: logger}, nil
}

func (l *zapLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.Debug(msg, mapToZapFields(fields)...)
}

func (l *zapLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.Info(msg, mapToZapFields(fields)...)
}

func (l *zapLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.Warn(msg, mapToZapFields(fields)...)
}

func (l *zapLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.Error(msg, mapToZapFields(fields)...)
}

func mapToZapFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return zapFields
}