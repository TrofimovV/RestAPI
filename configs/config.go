package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type ConfigDatabase struct {
	Port     string `env:"DB_PORT" env-default:"5432"`
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Username string `env:"DB_USER_NAME" env-default:"user"`
	Name     string `env:"DB_NAME" env-default:"postgres"`
	Password string `env:"DB_PASSWORD"`
}

var cfg ConfigDatabase

//Read config from .env file

func NewConfig() *ConfigDatabase {
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Print(err)
	}
	return &cfg
}
