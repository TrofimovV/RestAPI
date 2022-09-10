package configs

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type ConfigDatabase struct {
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Username string `env:"POSTGRES_USER" env-default:"user"`
	DBName   string `env:"POSTGRES_DB" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD"`
}

var cfg ConfigDatabase

//Read config from .env file

func NewConfig() *ConfigDatabase {
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		log.Print(err)
	}
	return &cfg
}
