package config

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env      string `envconfig:"ENV" default:"local"`
	Auth     Auth
	Todolist Todolist
	App      App
}

type Auth struct {
	Addr string `envconfig:"AUTH_ADDR" required:"true"`
}
type Todolist struct {
	Addr string `envconfig:"TODOLIST_ADDR" required:"true"`
}

type App struct {
	Addr          string `envconfig:"APP_ADDR" required:"true"`
	PublicKeyPath string `envconfig:"APP_PUBLIC_KEY_PATH" required:"true"`
}

func MustNew() *Config {
	cfg := new(Config)

	os.Setenv("ENV", "local")
	os.Setenv("AUTH_ADDR", ":8081")
	os.Setenv("TODOLIST_ADDR", ":8082")
	os.Setenv("APP_ADDR", "8080")
	os.Setenv("APP_PUBLIC_KEY_PATH", "./../cert/public.pem")

	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("config error: %v\n", err)
	}
	log.Printf("config: %+v", cfg)

	return cfg

}
