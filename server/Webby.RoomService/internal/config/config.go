package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string      `yaml:"env" env-default:"local"`
	ConnectionString string      `yaml:"connectionString"`
	JwtSecret        string      `yaml:"jwtSecret"`
	Http             HttpConfig  `yaml:"http"`
	Aws              AwsConfig   `yaml:"aws"`
	Grpc             GrpcConfig  `yaml:"grpc"`
	Redis            RedisConfig `yaml:"redis"`
}

type HttpConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type AwsConfig struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
}

type GrpcConfig struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	MediaServiceAddress    string `yaml:"mediaServiceAddress"`
	ChatServiceAddress     string `yaml:"chatServiceAddress"`
	CategoryServiceAddress string `yaml:"categoryServiceAddress"`
	QueueServiceAddress    string `yaml:"queueServiceAddress"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" env-default:"localhost:6379"`
	Password string `yaml:"password" env-default:""`
	DB       int    `yaml:"db" env-default:"0"`
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
