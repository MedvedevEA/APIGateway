package main

import (
	"APIGetway/internal/config"
	"APIGetway/internal/logger"
)

func main() {
	cfg := config.MustNew()

	lg := logger.MustNew(cfg.Env)

}
