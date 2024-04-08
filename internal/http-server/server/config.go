package server

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env        string `yaml:"env"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address"`
}

func NewConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("Z:/Golang/New_API/configs/configs.yaml", &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
