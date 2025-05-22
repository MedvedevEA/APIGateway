package main

import (
	"ppApiGatewayService/internal/config"
	"ppApiGatewayService/internal/logger"
)

func main() {
	cfg := config.MustNew()

	lg := logger.MustNew(cfg.Env)

}
