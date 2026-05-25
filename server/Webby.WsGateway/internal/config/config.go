package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string        `yaml:"env" env-default:"local"`
	JwtSecret        string        `yaml:"jwtSecret" env:"JWT_SECRET" env-required:"true"`
	ConnectionString string        `yaml:"connectionString"`
	TokenTTL         time.Duration `yaml:"token_ttl" env-default:"15m"`

	Http struct {
		Host    string        `yaml:"host" env-default:"localhost"`
		Port    int           `yaml:"port" env-default:"8090"`
		Timeout time.Duration `yaml:"timeout" env-default:"10s"`
	} `yaml:"http"`

	Redis struct {
		Addr     string `yaml:"addr" env-default:"localhost:6379"`
		Password string `yaml:"password" env-default:""`
		DB       int    `yaml:"db" env-default:"0"`
		Pattern  string `yaml:"pattern" env-default:"room:*"`
	} `yaml:"redis"`

	Grpc struct {
		Chat  string `yaml:"chat"`
		Votes string `yaml:"votes"`
		Room  string `yaml:"room"`
	} `yaml:"grpc"`

	Worker struct {
		ActivityInterval time.Duration `yaml:"activityInterval" env-default:"60s"`
	} `yaml:"worker"`
}

func MustLoad() *Config {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.yaml"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
