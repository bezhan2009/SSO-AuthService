package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env_default:"local"`
	StoragePath string        `yaml:"storage_path" env_required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env_default:"1h"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env_required:"true"`
	Timeout time.Duration `yaml:"timeout" env_required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(storagePath string) *Config {
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		panic("config file not exist:" + storagePath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(storagePath, &cfg); err != nil {
		panic("load config fail:" + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
