package config

import (
	"fmt"
	//	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"path/filepath"
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

func MustLoadFromPath(configPath string) *Config {

	return ReadConfig(configPath)

}

func MustLoad() *Config {

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		var config *Config
		configPath = "config/config_local.yaml"
		defer func() {
			if err := recover(); err != nil {
				configPath = "../config/config_local.yaml"
				config = ReadConfig(configPath)
			}
		}()
		config = ReadConfig(configPath)
		return config
	}

	return ReadConfig(configPath)

}

func ReadConfig(configPath string) *Config {

	if configPath == "" {
		panic("config path is empty")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current working directory: ", err)
	}

	absPath := filepath.Join(wd, configPath)

	_, err = os.Stat(configPath)

	if os.IsNotExist(err) {
		panic(fmt.Sprintf("config does not exist at path: %s (cwd: %s)", absPath, wd))
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		panic(fmt.Sprintf("failed to read config: %v (full path: %s)", err, absPath))
	}

	return &cfg
}
