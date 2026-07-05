package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string        `yaml:"env" env-default:"local" env:"ENV"`
	JwtSecret        string        `yaml:"jwtSecret" env:"JWT_SECRET" env-required:"true"`
	ConnectionString string        `yaml:"connectionString" env:"CONNECTION_STRING"`
	TokenTTL         time.Duration `yaml:"token_ttl" env-default:"15m" env:"TOKEN_TTL"`
	Http             Http          `yaml:"http"`
	Redis            Redis         `yaml:"redis"`
	Grpc             Grpc          `yaml:"grpc"`
}

type Http struct {
	HostPort    string        `yaml:"serviceAddress" env-default:"localhost:8090" env:"HTTP_WSGATEWAY_ADDRESS"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s" env:"HTTP_TIMEOUT"`
	CallTimeout time.Duration `yaml:"call_timeout" env-default:"5s" env:"HTTP_CALL_TIMEOUT"`
}

type Redis struct {
	Addr     string `yaml:"addr" env-default:"localhost:6379" env:"REDIS_ADDR"`
	Password string `yaml:"password" env-default:"" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env-default:"0" env:"REDIS_DB"`
	Pattern  string `yaml:"pattern" env-default:"*:*" env:"REDIS_PATTERN"`
}

type Grpc struct {
	ChatServiceAddress string `yaml:"chatServiceAddress" env-default:"localhost:50051" env:"GRPC_CHAT_SERVICE_ADDRESS"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()

	if configPath == "" {
		var cfg Config
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			panic("cannot read config from env: " + err.Error())
		}
		return &cfg
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
