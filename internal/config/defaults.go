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
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "text",
			Output:     "both",
			FilePath:   "./logs/clipd.log",
			MaxSize:    10, // 10MB per file
			MaxBackups: 3,  // Keep 3 old files
			MaxAge:     30, // Keep logs for 30 days
		},
	}
}
