package server

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"ppApiGatewayService/internal/config"
	"ppApiGatewayService/internal/server/middleware"
	"ppApiGatewayService/internal/service"
	"ppApiGatewayService/pkg/secure"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	app *fiber.App
	lg  *slog.Logger
	cfg *config.Server
}

func MustNew(svc *service.Service, lg *slog.Logger, cfg *config.Server) *Server {
	publicKey, err := secure.LoadPublicKey(cfg.PublicKeyPath)
	if err != nil {
		log.Fatalf("failed to initialize server: %v\n", err)
	}
	app := fiber.New(fiber.Config{
		AppName:      cfg.Name,
		WriteTimeout: cfg.WriteTimeout,
	})
	app.Use(recover.New(recover.ConfigDefault))
	app.Use(middleware.GetLoggerMiddlewareFunc(lg))

	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")
	//public
	v1Group.Post("/login", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	v1Group.Post("/registration", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	//refresh
	v1Group.Post(
		"/refresh",
		middleware.GetAuthMiddlewareFunc(publicKey, "refresh"),
		func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) },
	)
	//access
	authGroup := v1Group.Group(
		"/auth",
		middleware.GetAuthMiddlewareFunc(publicKey, "access"),
	)

	authGroup.Post("/logout", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })
	authGroup.Post("/unregistration", func(ctx *fiber.Ctx) error { return ctx.SendStatus(fiber.StatusNotImplemented) })

	authGroup.Post("/tasks", svc.AddTask)
	authGroup.Get("/tasks/:taskId", svc.GetTask)
	authGroup.Get("/tasks", svc.GetTasks)
	authGroup.Patch("/tasks/:taskId", svc.UpdateTask)
	authGroup.Delete("/tasks/:taskId", svc.RemoveTask)

	app.Use(middleware.BadRequest)

	return &Server{
		app,
		lg,
		cfg,
	}
}

func (s *Server) Start() {
	chErr := make(chan error, 1)
	defer close(chErr)
	go func() {
		s.lg.Info("server is started", slog.String("owner", "server"), slog.String("bindAddress", s.cfg.BindAddr))
		if err := s.app.Listen(s.cfg.BindAddr); err != nil {
			chErr <- err
		}
	}()
	go func() {
		chQuit := make(chan os.Signal, 1)
		signal.Notify(chQuit, syscall.SIGINT, syscall.SIGTERM)
		<-chQuit
		chErr <- s.app.Shutdown()
	}()
	if err := <-chErr; err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "server"))
		return
	}
	s.lg.Info("server is stoped", slog.String("owner", "server"))

}
