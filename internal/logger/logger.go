package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/geodask/clipboard-manager/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg config.LoggingConfig) (*slog.Logger, error) {
	var logLevel slog.Level

	switch cfg.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var writers []io.Writer

	if cfg.Output == "file" || cfg.Output == "both" {
		logDir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}

		fileWriter := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,    // megabytes
			MaxBackups: cfg.MaxBackups, // number of old files
			MaxAge:     cfg.MaxAge,     // days
			Compress:   true,           // compress old logs
		}
		writers = append(writers, fileWriter)
	}

	if cfg.Output == "stdout" || cfg.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	writer := io.MultiWriter(writers...)

	var handler slog.Handler
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(writer, opts)
	case "text":
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewTextHandler(writer, opts)
	}

	return slog.New(handler), nil
}
