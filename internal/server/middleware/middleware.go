package middleware

import (
	"crypto/rsa"
	"log/slog"

	srvErr "ppApiGatewayService/internal/server/err"
	"ppApiGatewayService/pkg/jwt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetLoggerMiddlewareFunc(lg *slog.Logger) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		err := ctx.Next()

		lg.Info(
			"incoming server request",
			slog.String("owner", "server"),
			slog.Any("method", ctx.Method()),
			slog.Any("path", ctx.Path()),
			slog.Any("statusCode", ctx.Response().StatusCode()),
			slog.Any("time", time.Since(start)),
		)
		return err
	}
}
func GetAuthMiddlewareFunc(publicKey *rsa.PublicKey, tokenType string) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get("Authorization")
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			return ctx.Status(401).SendString(srvErr.ErrInvalidTokenFormat.Error())
		}
		tokenClaims, err := jwt.ParseToken(tokenString[7:], publicKey)
		if err != nil {
			return ctx.Status(401).SendString(srvErr.ErrInvalidToken.Error())
		}
		if tokenClaims.TokenType != tokenType {
			return ctx.Status(401).SendString(srvErr.ErrInvalidTokenType.Error())
		}
		ctx.Locals("claims", tokenClaims)
		return ctx.Next()
	}
}
func BadRequest(ctx *fiber.Ctx) error {
	return ctx.Status(404).SendString(srvErr.ErrRouteNotFound.Error())
}
