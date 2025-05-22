package service

import (
	"log/slog"
	"ppApiGatewayService/internal/todolist"

	proto "github.com/MedvedevEA/ppProtos/gen/auth"
)

type Service struct {
	auth     proto.AuthServiceClient
	todolist todolist.Repository
	lg       *slog.Logger
}

func MustNew(auth proto.AuthServiceClient, todolist todolist.Repository, lg *slog.Logger) *Service {
	return &Service{
		auth,
		todolist,
		lg,
	}
}
