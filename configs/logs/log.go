package logs

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutputLogs()},
		Level:       zap.NewAtomicLevelAt(getLogLevel()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	log, err = logConfig.Build()
	if err != nil {
		panic(err)
	}
}

func Info(message string, tags ...zap.Field) {
	log.Info(message, tags...)
	_ = log.Sync()
}

func Debug(message string, tags ...zap.Field) {
	log.Debug(message, tags...)
	_ = log.Sync()
}

func Warn(message string, tags ...zap.Field) {
	log.Warn(message, tags...)
	_ = log.Sync()
}

func Error(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Error(message, tags...)
	_ = log.Sync()
}

func Fatal(message string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Fatal(message, tags...)
	_ = log.Sync()
}

func getOutputLogs() string {
	output := strings.TrimSpace(os.Getenv("LOG_OUTPUT"))
	if output == "" {
		return "stdout"
	}

	return strings.ToLower(output)
}

func getLogLevel() zapcore.Level {
	level := strings.TrimSpace(os.Getenv("LOG_LEVEL"))

	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// SetLogLevel dynamically changes the log level at runtime
func SetLogLevel(level string) {
	logLevel := getLogLevelFromString(level)
	logConfig := zap.NewAtomicLevelAt(logLevel)
	logConfig.SetLevel(logLevel)
}

func getLogLevelFromString(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
