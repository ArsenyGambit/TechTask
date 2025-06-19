package config

import "time"

type Config struct {
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
		Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
		User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-default:"password"`
		DBName   string `yaml:"dbname" env:"DB_NAME" env-default:"news_db"`
		SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
	} `yaml:"database"`

	Server struct {
		GRPCPort int `yaml:"grpc_port" env:"GRPC_PORT" env-default:"8080"`
	} `yaml:"server"`

	Cache struct {
		TTL time.Duration `yaml:"ttl" env:"CACHE_TTL" env-default:"5m"`
	} `yaml:"cache"`
}
