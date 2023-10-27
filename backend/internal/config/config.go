package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port         string `env:"PORT" env-default:"3000"`
	RedisAddress string `env:"REDISADDR" env-default:"redis:6379"`
	LoggerLvl    string `env:"LOGGERLVL" env-default:"local"`
	DSN          string `env:"DSN" env-default:"secret"`
}

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
