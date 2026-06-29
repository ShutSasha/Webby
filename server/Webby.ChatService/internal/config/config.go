package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string      `yaml:"env" env-default:"local" env:"ENV"`
	ConnectionString string      `yaml:"connectionString" env:"CONNECTION_STRING"`
	JwtSecret        string      `yaml:"jwtSecret" env:"JWT_SECRET"`
	Http             HttpConfig  `yaml:"http"`
	Grpc             GrpcConfig  `yaml:"grpc"`
	Redis            RedisConfig `yaml:"redis"`
}

type HttpConfig struct {
	Host    string        `yaml:"host" env-default:"0.0.0.0" env:"HTTP_HOST"`
	Port    int           `yaml:"port" env-default:"8081" env:"HTTP_PORT"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s" env:"HTTP_TIMEOUT"`
}

type GrpcConfig struct {
	Port               int    `yaml:"port" env-default:"50051" env:"GRPC_PORT"`
	RoomServiceAddress string `yaml:"roomServiceAddress" env-default:"localhost:5006" env:"GRPC_ROOM_SERVICE_ADDRESS"`
	UserServiceAddress string `yaml:"userServiceAddress" env-default:"localhost:5004" env:"GRPC_USER_SERVICE_ADDRESS"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" env-default:"localhost:6379" env:"REDIS_ADDR"`
	Password string `yaml:"password" env-default:"" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env-default:"0" env:"REDIS_DB"`
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
