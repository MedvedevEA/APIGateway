package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env      string `envconfig:"ENV" default:"local"`
	Auth     Auth
	Server   Server
	Todolist Todolist
}

type Auth struct {
	BindAddr string `envconfig:"AUTH_BIND_ADDR" required:"true"`
}
type Server struct {
	BindAddr      string        `envconfig:"SERVER_BIND_ADDR" required:"true"`
	Name          string        `envconfig:"SERVER_NAME" required:"true"`
	WriteTimeout  time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" required:"true"`
	PublicKeyPath string        `envconfig:"SERVER_PUBLIC_KEY_PATH" required:"true"`
}
type Todolist struct {
	BindAddr string `envconfig:"TODOLIST_BIND_ADDR" required:"true"`
}
type Token struct {
}

func MustNew() *Config {
	//TODO
	if err := godotenv.Load("./../../.env"); err != nil {
		log.Fatalf("failed to load configuration: %v\n", err)
	}

	cfg := new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("failed to load configuration: %v\n", err)
	}
	return cfg
}
