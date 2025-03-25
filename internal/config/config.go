package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	AppParams      AppParams      `yaml:"app_params"`
	SqliteParams   SqliteParams   `yaml:"sqlite_params" env_required:"true"`
	PostgresParams PostgresParams `yaml:"postgres_params" env_required:"true"`
	RedisParams    RedisParams    `yaml:"redis_params" env_required:"true"`
	AuthParams     AuthParams     `yaml:"auth_params" env_required:"true"`
	GRPC           GRPCConfig     `yaml:"grpc"`
}

type AppParams struct {
	Env  string `yaml:"env" env_default:"local"`
	DBSM string `yaml:"dbsm" env_default:"postgres"`
}

type PostgresParams struct {
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type RedisParams struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type SqliteParams struct {
	StoragePath string `yaml:"storage_path"`
}

type AuthParams struct {
	JwtTTLMinutes      time.Duration `yaml:"jwt_ttl_minutes"`
	JwtTTLRefreshHours time.Duration `yaml:"jwt_ttl_refresh_hours"`
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
