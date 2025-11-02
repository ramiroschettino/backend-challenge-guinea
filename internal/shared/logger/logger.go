package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger es la interfaz que usamos en toda la aplicación
// Nos permite cambiar la implementación sin tocar el resto del código
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}

// zapLogger es la implementación concreta usando Zap (librería de Uber)
type zapLogger struct {
	logger *zap.Logger
}

// NewLogger crea un nuevo logger configurado
// level: debug, info, warn, error
// format: json (para producción) o text (para desarrollo)
func NewLogger(level, format string) (Logger, error) {
	var config zap.Config

	// Configuración según formato
	if format == "json" {
		config = zap.NewProductionConfig() // Logs estructurados en JSON
	} else {
		config = zap.NewDevelopmentConfig() // Logs legibles para humanos
	}

	// Configurar nivel de log
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

// Implementación de los métodos de la interfaz Logger
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

// mapToZapFields convierte un map a campos de Zap
// Ejemplo: {"user_id": "123", "action": "create"} → [Field, Field]
func mapToZapFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return zapFields
}