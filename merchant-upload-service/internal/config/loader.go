package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*AppConfig, error) {
	if path == "" {
		if cfg, err := LoadFromEnv(); err == nil {
			return cfg, nil
		}

		defaultPaths := []string{
			"config.yaml",
			"configs/config.yaml",
			"../config.yaml",
			"../configs/config.yaml",
		}
		for _, p := range defaultPaths {
			if cfg, err := LoadFromFile(p); err == nil {
				return cfg, nil
			}
		}
		return nil, errors.New("load config file failed")
	}
	return LoadFromFile(path)
}

func LoadFromFile(path string) (*AppConfig, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("file not exist")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %s failed: %w", path, err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config file %s failed: %w", path, err)
	}

	if err := ValidateServerConfig(&cfg.Server); err != nil {
		return nil, fmt.Errorf("validate server config failed: %w", err)
	}

	if err := ValidateDBConfig(&cfg.Database); err != nil {
		return nil, fmt.Errorf("validate db config failed: %w", err)
	}

	return &cfg, nil
}

func LoadFromEnv() (*AppConfig, error) {
	cfg := &AppConfig{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "localhost"),
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  Duration(getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second)),
			WriteTimeout: Duration(getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second)),
		},
		Database: DBConfig{
			Dialect:  getEnv("DB_DIALECT", "sqlite"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 3306),
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", "123456"),
			Database: getEnv("DB_DATABASE", "merchant"),

			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: Duration(getEnvDuration("DB_CONN_MAX_LIFETIME", 1*time.Hour)),
			ConnMaxIdleTime: Duration(getEnvDuration("DB_CONN_MAX_IDLE_TIME", 30*time.Minute)),
		},
	}

	if err := ValidateConfig(cfg); err != nil {
		return nil, fmt.Errorf("validate config failed: %w", err)
	}
	return cfg, nil
}

func ValidateConfig(cfg *AppConfig) error {
	if err := ValidateServerConfig(&cfg.Server); err != nil {
		return fmt.Errorf("validate server config failed: %w", err)
	}

	if err := ValidateDBConfig(&cfg.Database); err != nil {
		return fmt.Errorf("validate db config failed: %w", err)
	}
	return nil
}

func ValidateServerConfig(cfg *ServerConfig) error {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port %d", cfg.Port)
	}

	if cfg.ReadTimeout <= 0 {
		return fmt.Errorf("invalid read timeout %d", cfg.ReadTimeout)
	}

	if cfg.WriteTimeout <= 0 {
		return fmt.Errorf("invalid write timeout %d", cfg.WriteTimeout)
	}
	return nil
}

func ValidateDBConfig(cfg *DBConfig) error {
	if cfg.Dialect == "" {
		return fmt.Errorf("dialect is empty")
	}
	if cfg.Host == "" {
		return fmt.Errorf("host is empty")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port %d", cfg.Port)
	}
	if cfg.Username == "" {
		return fmt.Errorf("username is empty")
	}
	if cfg.Password == "" {
		return fmt.Errorf("password is empty")
	}
	if cfg.Database == "" {
		return fmt.Errorf("database is empty")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var i int
		if err := yaml.Unmarshal([]byte(value), &i); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		var d time.Duration
		if err := yaml.Unmarshal([]byte(value), &d); err == nil {
			return d
		}
	}
	return defaultValue
}
