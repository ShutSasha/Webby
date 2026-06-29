package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string       `yaml:"env" env-default:"local" env:"ENV"`
	ConnectionString string       `yaml:"connectionString" env:"CONNECTION_STRING"`
	JwtSecret        string       `yaml:"jwtSecret" env:"JWT_SECRET"`
	Worker           WorkerConfig `yaml:"worker"`
	Http             HttpConfig   `yaml:"http"`
	Aws              AwsConfig    `yaml:"aws"`
	Grpc             GrpcConfig   `yaml:"grpc"`
	Redis            RedisConfig  `yaml:"redis"`
}

type WorkerConfig struct {
	PointsPerTick int           `yaml:"points_per_tick" env-default:"15" env:"POINTS_PER_TICK"`
	Interval      time.Duration `yaml:"interval" env-default:"1m" env:"INTERVAL"`
	ZombieTTL     time.Duration `yaml:"zombie_ttl" env-default:"2m" env:"ZOMBIE_TTL"`
}
type HttpConfig struct {
	Host    string        `yaml:"host" env-default:"0.0.0.0" env:"HTTP_HOST"`
	Port    int           `yaml:"port" env-default:"8080" env:"HTTP_PORT"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s" env:"HTTP_TIMEOUT"`
}

type AwsConfig struct {
	AccessKey string `yaml:"accessKey" env:"AWS_ACCESS_KEY"`
	SecretKey string `yaml:"secretKey" env:"AWS_SECRET_KEY"`
	Region    string `yaml:"region" env:"AWS_REGION"`
	Bucket    string `yaml:"bucket" env:"AWS_BUCKET"`
}

type GrpcConfig struct {
	Host                       string `yaml:"host" env-default:"localhost" env:"GRPC_HOST"`
	Port                       int    `yaml:"port" env-default:"50054" env:"GRPC_PORT"`
	MediaServiceAddress        string `yaml:"mediaServiceAddress" env-default:"localhost:5005" env:"GRPC_MEDIA_SERVICE_ADDRESS"`
	ChatServiceAddress         string `yaml:"chatServiceAddress" env-default:"localhost:5008" env:"GRPC_CHAT_SERVICE_ADDRESS"`
	CategoryServiceAddress     string `yaml:"categoryServiceAddress" env-default:"localhost:5009" env:"GRPC_CATEGORY_SERVICE_ADDRESS"`
	QueueServiceAddress        string `yaml:"queueServiceAddress" env-default:"localhost:5007" env:"GRPC_QUEUE_SERVICE_ADDRESS"`
	NotificationServiceAddress string `yaml:"notificationServiceAddress" env-default:"localhost:5007" env:"GRPC_NOTIFICATION_SERVICE_ADDRESS"`
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
