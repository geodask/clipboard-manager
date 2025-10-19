package config

import "time"

func Default() *Config {
	return &Config{
		Database: DatabaseConfig{
			Path: "./clipboard.db",
		},
		API: APIConfig{
			SocketPath:   "/tmp/clipd.sock",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
		Monitor: MonitorConfig{},
		Daemon: DaemonConfig{
			ShutdownTimeout: 5 * time.Second,
			PollInterval:    500 * time.Millisecond,
		},
	}
}
