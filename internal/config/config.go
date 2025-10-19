package config

import (
	"flag"
	"time"
)

type Config struct {
	Database DatabaseConfig
	API      APIConfig
	Monitor  MonitorConfig
	Daemon   DaemonConfig
	Logging  LoggingConfig
}

type DatabaseConfig struct {
	Path string
}

type APIConfig struct {
	SocketPath   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type MonitorConfig struct {
}

type DaemonConfig struct {
	ShutdownTimeout time.Duration
	PollInterval    time.Duration
}

type LoggingConfig struct {
	Level      string
	Format     string
	Output     string
	FilePath   string
	MaxSize    int // megabytes
	MaxBackups int // number of old files to keep
	MaxAge     int // days
}

func Load() (*Config, error) {
	cfg := Default()

	flag.StringVar(&cfg.Database.Path, "db", cfg.Database.Path, "Path to SQLite database")

	flag.StringVar(&cfg.API.SocketPath, "socket", cfg.API.SocketPath, "Path to Unix socket for API")
	flag.DurationVar(&cfg.API.ReadTimeout, "read-timeout", cfg.API.ReadTimeout, "HTTP read timeout")
	flag.DurationVar(&cfg.API.WriteTimeout, "write-timeout", cfg.API.WriteTimeout, "HTTP write timeout")
	flag.DurationVar(&cfg.API.IdleTimeout, "idle-timeout", cfg.API.IdleTimeout, "HTTP idle timeout")

	flag.DurationVar(&cfg.Daemon.PollInterval, "poll-interval", cfg.Daemon.PollInterval, "Clipboard polling interval")
	flag.DurationVar(&cfg.Daemon.ShutdownTimeout, "shutdown-timeout", cfg.Daemon.ShutdownTimeout, "Graceful shutdown timeout")

	flag.StringVar(&cfg.Logging.Level, "log-level", cfg.Logging.Level, "Log level (debug, info, warn, error)")
	flag.StringVar(&cfg.Logging.Format, "log-format", cfg.Logging.Format, "Log format (text, json)")
	flag.StringVar(&cfg.Logging.Output, "log-output", cfg.Logging.Output, "Log output (stdout, file, both)")
	flag.StringVar(&cfg.Logging.FilePath, "log-file", cfg.Logging.FilePath, "Log file path")
	flag.IntVar(&cfg.Logging.MaxSize, "log-max-size", cfg.Logging.MaxSize, "Max log file size in MB")
	flag.IntVar(&cfg.Logging.MaxBackups, "log-max-backups", cfg.Logging.MaxBackups, "Max number of old log files")
	flag.IntVar(&cfg.Logging.MaxAge, "log-max-age", cfg.Logging.MaxAge, "Max age of log files in days")
	flag.Parse()

	return cfg, nil
}
