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

func Load() (*Config, error) {
	cfg := Default()

	flag.StringVar(&cfg.Database.Path, "db", cfg.Database.Path, "Path to SQLite database")
	flag.StringVar(&cfg.API.SocketPath, "socket", cfg.API.SocketPath, "Path to Unix socket for API")
	flag.DurationVar(&cfg.API.ReadTimeout, "read-timeout", cfg.API.ReadTimeout, "HTTP read timeout")
	flag.DurationVar(&cfg.API.WriteTimeout, "write-timeout", cfg.API.WriteTimeout, "HTTP write timeout")
	flag.DurationVar(&cfg.API.IdleTimeout, "idle-timeout", cfg.API.IdleTimeout, "HTTP idle timeout")
	flag.DurationVar(&cfg.Daemon.PollInterval, "poll-interval", cfg.Daemon.PollInterval, "Clipboard polling interval")
	flag.DurationVar(&cfg.Daemon.ShutdownTimeout, "shutdown-timeout", cfg.Daemon.ShutdownTimeout, "Graceful shutdown timeout")

	flag.Parse()

	return cfg, nil
}
