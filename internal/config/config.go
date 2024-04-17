package config

import (
	//	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env           string           `yaml:"env" env-default:"local"`
	GRPC          GRPCConfig       `yaml:"grpc"`
	MigrationPath string           `yaml:"migrations_path"`
	TokenTTL      time.Duration    `yaml:"token_ttl" env-default:"1h"`
	Postgres      PostgreSQLConfig `yaml:"postgres" env-required:"true"`
}
type PostgreSQLConfig struct {
	Host     string `yaml:"host" env:"PGHOST" env-default:"0.0.0.0"`
	Port     int    `yaml:"port" env:"PGPORT" env-default:"5555"`
	User     string `yaml:"user" env:"PGUSER" env-default:"postgres"`
	Password string `yaml:"password" env:"PGPASSWORD"`
	DBName   string `yaml:"dbname" env:"PGDATABASE"`
	SSLMode  string `yaml:"sslmode" env:"PGSSLMODE" env-default:"disable"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout" env-default:"1h"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		panic("config path empty")
	}

	_, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		panic("config does not exist in " + configPath)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	return configPath
}
