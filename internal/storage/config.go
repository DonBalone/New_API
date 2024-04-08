package storage

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	StorageInfo `yaml:"storage_info"`
}

type StorageInfo struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	SSLMode  string `yaml:"sslmode"`
	DBName   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func NewConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("Z:/Golang/New_API/configs/configs.yaml", &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
