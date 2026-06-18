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
	Host    string        `yaml:"host" env-default:"localhost" env:"HTTP_HOST"`
	Port    int           `yaml:"port" env-default:"5050" env:"HTTP_PORT"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s" env:"HTTP_TIMEOUT"`
}

type GrpcConfig struct {
	Host                       string `yaml:"host" env-default:"localhost" env:"GRPC_HOST"`
	Port                       int    `yaml:"port" env-default:"5051" env:"GRPC_PORT"`
	MediaServiceAddress        string `yaml:"mediaServiceAddress" env-default:"localhost:5005" env:"GRPC_MEDIA_SERVICE_ADDRESS"`
	ComplaintServiceAddress    string `yaml:"complaintServiceAddress" env-default:"localhost:5004" env:"GRPC_COMPLAINT_SERVICE_ADDRESS"`
	NotificationServiceAddress string `yaml:"notificationServiceAddress" env-default:"localhost:5007" env:"GRPC_NOTIFICATION_SERVICE_ADDRESS"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
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
