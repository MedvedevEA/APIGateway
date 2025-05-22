package service

import (
	"log/slog"
	repoAuth "ppApiGatewayService/internal/repository/auth"
	repoTodolist "ppApiGatewayService/internal/repository/todolist"
)

type Service struct {
	auth     repoAuth.Repository
	todolist repoTodolist.Repository
	lg       *slog.Logger
}

func MustNew(auth repoAuth.Repository, todolist repoTodolist.Repository, lg *slog.Logger) *Service {
	return &Service{
		auth,
		todolist,
		lg,
	}
}
