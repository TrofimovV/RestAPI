package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Host     string `env:"HOST"`
	Username string `env:"USER_NAME"`
	DB       string `env:"DB_NAME"`
	Password string `env:"DB_PASSWORD"'`
}

var cfg *Config

func NewConfig() *Config {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Print(err)
	}
	return cfg
}
