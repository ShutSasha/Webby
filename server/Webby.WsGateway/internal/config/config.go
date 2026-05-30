package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string        `yaml:"env" env-default:"local" env:"ENV"`
	JwtSecret        string        `yaml:"jwtSecret" env:"JWT_SECRET" env-required:"true"`
	ConnectionString string        `yaml:"connectionString" env:"CONNECTION_STRING"`
	TokenTTL         time.Duration `yaml:"token_ttl" env-default:"15m" env:"TOKEN_TTL"`

	Http struct {
		Host    string        `yaml:"host" env-default:"localhost" env:"HTTP_HOST"`
		Port    int           `yaml:"port" env-default:"8090" env:"HTTP_PORT"`
		Timeout time.Duration `yaml:"timeout" env-default:"10s" env:"HTTP_TIMEOUT"`
	} `yaml:"http"`

	Redis struct {
		Addr     string `yaml:"addr" env-default:"localhost:6379" env:"REDIS_ADDR"`
		Password string `yaml:"password" env-default:"" env:"REDIS_PASSWORD"`
		DB       int    `yaml:"db" env-default:"0" env:"REDIS_DB"`
		Pattern  string `yaml:"pattern" env-default:"room:*" env:"REDIS_PATTERN"`
	} `yaml:"redis"`

	Grpc struct {
		Chat  string `yaml:"chat" env:"GRPC_CHAT"`
		Votes string `yaml:"votes" env:"GRPC_VOTES"`
		Room  string `yaml:"room" env:"GRPC_ROOM"`
	} `yaml:"grpc"`
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
