package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host                string `yaml:"host"`
	Port                int    `yaml:"port"`
	JWTSecret           string `yaml:"jwt_secret"`
	JWTTTLHours         int    `yaml:"jwt_ttl_hours"`
	EncryptionKey       string `yaml:"encryption_key"`
	ReadTimeoutSeconds  int    `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `yaml:"write_timeout_seconds"`
}

type DatabaseConfig struct {
	Host                     string `yaml:"host"`
	Port                     int    `yaml:"port"`
	User                     string `yaml:"user"`
	Password                 string `yaml:"password"`
	Name                     string `yaml:"name"`
	SSLMode                  string `yaml:"sslmode"`
	MaxConns                 int    `yaml:"max_conns"`
	ConnectRetries           int    `yaml:"connect_retries"`
	ConnectRetryDelaySeconds int    `yaml:"connect_retry_delay_seconds"`
}

type DockerConfig struct {
	Enabled     bool   `yaml:"enabled"`
	ComposeFile string `yaml:"compose_file"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Docker   DockerConfig   `yaml:"docker"`
	CORS     CORSConfig     `yaml:"cors"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if len(cfg.Server.EncryptionKey) != 32 {
		return nil, fmt.Errorf("encryption_key must be exactly 32 bytes for AES-256, got %d", len(cfg.Server.EncryptionKey))
	}
	return &cfg, nil
}

func (c *Config) ReadTimeout() time.Duration {
	return time.Duration(c.Server.ReadTimeoutSeconds) * time.Second
}

func (c *Config) WriteTimeout() time.Duration {
	return time.Duration(c.Server.WriteTimeoutSeconds) * time.Second
}

func (c *Config) JWTTTL() time.Duration {
	return time.Duration(c.Server.JWTTTLHours) * time.Hour
}
