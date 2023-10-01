package config

import (
	"log"
	"sync"

	rediscache "github.com/ArturKKK/CurrencyTransfer/internal/cache/redis"
	"github.com/ArturKKK/CurrencyTransfer/internal/db"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"listen"`
	Postgres db.Config         `yaml:"postgres"`
	Redis    rediscache.Config `yaml:"redis"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			log.Fatal(err)
		}
	})

	return instance
}
