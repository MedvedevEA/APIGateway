package service

import (
	"context"
	"log/slog"
	svcDto "ppApiGatewayService/internal/service/dto"
	svcErr "ppApiGatewayService/internal/service/err"
	todoDto "ppApiGatewayService/internal/todolist/dto"
	"ppApiGatewayService/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/MedvedevEA/ppProtos/gen/auth"
)

func (s *Service) Registration(ctx *fiber.Ctx) error {
	req := new(svcDto.RegistrationRequest)
	if err := ctx.BodyParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Registration"))
		return ctx.Status(400).SendString(svcErr.ErrBodyParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RemoveMessage"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	authResp, err := s.auth.Register(
		context.Background(),
		&proto.RegisterRequest{
			Login:    req.Login,
			Password: req.Password,
		},
	)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Registration"))
		return ctx.Status(500).SendString(svcErr.ErrInternalServerError.Error())
	}

	var userId uuid.UUID
	if userId, err = uuid.Parse(authResp.UserId); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Registration"))
		return ctx.Status(500).SendString(svcErr.ErrInternalServerError.Error())
	}

	if err := s.todolist.AddUserWithUserId(&todoDto.AddUser{
		UserId: &userId,
		Name:   req.Name,
	}); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Registration"))
		return ctx.Status(500).SendString(svcErr.ErrInternalServerError.Error())
	}
	return ctx.SendStatus(204)
}
func (s *Service) Login(ctx *fiber.Ctx) error {
	req := new(svcDto.LoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Login"))
		return ctx.Status(400).SendString(svcErr.ErrBodyParse.Error())
	}
	if err := validator.Validate(req); err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Login"))
		return ctx.Status(400).SendString(svcErr.ErrValidate.Error())
	}
	authResp, err := s.auth.Login(
		context.Background(),
		&proto.LoginRequest{
			Login:      req.Login,
			Password:   req.Password,
			DeviceCode: req.DeviceCode,
		},
	)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Login"))
		st, ok := status.FromError(err)
		if !ok {
			return ctx.Status(500).SendString(svcErr.ErrInternalServerError.Error())
		}
		switch st.Code() {
		case codes.Unauthenticated:
			return ctx.Status(401).SendString(st.Message())
		case codes.InvalidArgument:
			return ctx.Status(400).SendString(st.Message())
		default:
			return ctx.Status(500).SendString(st.Message())
		}
	}
	return ctx.Status(200).JSON(&svcDto.LoginResponse{
		AccessToken:  authResp.AccessToken,
		RefreshToken: authResp.RefreshToken,
	})
}
