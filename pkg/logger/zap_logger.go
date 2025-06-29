package logger

import (
	"gptBot/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	log *zap.Logger
}

func NewLogger(cfg config.LoggerConfig) (*ZapLogger, error) {
	var zapCfg zap.Config
	if cfg.Environment == "production" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	zapCfg.Level = zap.NewAtomicLevelAt(parseLevel(cfg.Level))

	z, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{log: z}, nil
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
	l.log.Info(msg, convertFields(fields)...)
}

func (l *ZapLogger) Debug(msg string, fields ...Field) {
	l.log.Debug(msg, convertFields(fields)...)
}

func (l *ZapLogger) Warn(msg string, fields ...Field) {
	l.log.Warn(msg, convertFields(fields)...)
}

func (l *ZapLogger) Error(msg string, fields ...Field) {
	l.log.Error(msg, convertFields(fields)...)
}

func (l *ZapLogger) With(fields ...Field) Logger {
	child := l.log.With(convertFields(fields)...)
	return &ZapLogger{log: child}
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func convertFields(fields []Field) []zap.Field {
	zFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zFields[i] = zap.Any(f.Key, f.Value)
	}
	return zFields
}