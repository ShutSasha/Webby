package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string     `yaml:"env" env-default:"local" env:"ENV"`
	ConnectionString string     `yaml:"connectionString" env:"CONNECTION_STRING"`
	JwtSecret        string     `yaml:"jwtSecret" env:"JWT_SECRET"`
	Http             HttpConfig `yaml:"http"`
	Grpc             GrpcConfig `yaml:"grpc"`
}

type HttpConfig struct {
	HostPort string        `yaml:"serviceAddress" env-default:"localhost:8082" env:"HTTP_ROOM_CATEGORY_ADDRESS"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s" env:"HTTP_TIMEOUT"`
}

type GrpcConfig struct {
	HostPort string `yaml:"roomCategoryServiceAddress" env-default:"localhost:50052" env:"GRPC_ROOM_CATEGORY_SERVICE_ADDRESS"`
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
