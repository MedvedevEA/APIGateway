package main

import (
	"log"
	"ppApiGatewayService/internal/config"
	"ppApiGatewayService/internal/logger"
	"ppApiGatewayService/internal/server"
	"ppApiGatewayService/internal/service"
	"ppApiGatewayService/internal/todolist"

	proto "github.com/MedvedevEA/ppProtos/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustNew()
	lg := logger.MustNew(cfg.Env)

	conn, err := grpc.NewClient(cfg.Auth.BindAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to initialize auth service: %v\n", err)
	}
	defer conn.Close()
	auth := proto.NewAuthServiceClient(conn)

	todolist := todolist.MustNew(lg, &cfg.Todolist)
	service := service.MustNew(auth, todolist, lg)
	server := server.MustNew(service, lg, &cfg.Server)

	server.Start()

}
