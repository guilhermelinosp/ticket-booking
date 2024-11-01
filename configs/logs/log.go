package logs

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

func init() {
	// Cria a pasta de logs se não existir
	createLogDir()

	// Configura o logger para registrar requisições e erros em um único arquivo
	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(getLogLevel()),
		Encoding:    "json",
		OutputPaths: []string{getCombinedLogFile()},
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

func createLogDir() {
	// Define o caminho do diretório de logs
	logDir := "logs"

	// Tenta criar o diretório de logs, se não existir
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic("failed to create logs directory: " + err.Error())
	}
}

func Sync() {
	_ = log.Sync()
}

func Request(method, path string, status int, duration string) {
	log.Info("HTTP Request",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", status),
		zap.String("duration", duration),
	)
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

func getCombinedLogFile() string {
	// Gera o nome do arquivo de log baseado na data atual
	date := time.Now().Format("20060102")     // Formato: YYYYMMDD
	return filepath.Join("logs", date+".log") // Define o caminho do log com apenas YYYYMMDD.log
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
