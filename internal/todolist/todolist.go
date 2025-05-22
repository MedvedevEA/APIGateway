package todolist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"ppApiGatewayService/internal/config"
	todoDto "ppApiGatewayService/internal/todolist/dto"
	todoErr "ppApiGatewayService/internal/todolist/err"

	"github.com/gofiber/fiber/v2"
)

type Repository interface {
	AddUserWithUserId(dto *todoDto.AddUser) error

	Redirect(ctx *fiber.Ctx, entity string) error
}

type Todolist struct {
	bindAddr string
	client   *http.Client
	lg       *slog.Logger
}

func MustNew(lg *slog.Logger, cfg *config.Todolist) *Todolist {
	client := http.Client{
		Timeout: cfg.WriteTimeout,
	}
	return &Todolist{
		bindAddr: cfg.BindAddr,
		client:   &client,
		lg:       lg,
	}
}

func (t *Todolist) request(method string, host string, path string, param map[string]string, query map[string]string, body any) (*http.Response, error) {

	reqUrl, err := url.Parse(host + "/" + path)
	if err != nil {
		return nil, err
	}
	for _, value := range param {
		reqUrl.Path = reqUrl.Path + "/" + value
	}

	if query != nil {
		urlValues := url.Values{}
		for key, value := range query {
			urlValues.Add(key, value)
		}
		reqUrl.RawQuery = urlValues.Encode()
	}
	var reqBody []byte
	if body != nil {
		if b, ok := body.([]byte); ok {
			reqBody = b
		} else {
			reqBody, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
		}
	}

	req, err := http.NewRequest(method, reqUrl.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return t.client.Do(req)

}

func (t *Todolist) AddUserWithUserId(dto *todoDto.AddUser) error {
	resp, err := t.request("POST", t.bindAddr, "users", nil, nil, dto)
	if err != nil {
		t.lg.Error(err.Error(), slog.String("owner", "todolist.AddUserWithUserId"))
		return todoErr.ErrInternalServerError
	}
	if resp.StatusCode != 201 {
		msg := fmt.Sprintf("want StatusCode=204 have StatusCode=%d", resp.StatusCode)
		t.lg.Error(msg, slog.String("owner", "todolist.AddUserWithUserId"))
		return todoErr.ErrInternalServerError
	}
	return nil
}

func (t *Todolist) Redirect(ctx *fiber.Ctx, entity string) error {
	resp, err := t.request(ctx.Method(), t.bindAddr, entity, ctx.AllParams(), ctx.Queries(), ctx.Body())
	if err != nil {
		t.lg.Error(err.Error(), slog.String("owner", "todolist.Redirect"))
		return ctx.Status(500).SendString(todoErr.ErrInternalServerError.Error())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.lg.Error(err.Error(), slog.String("owner", "todolist.Redirect"))
		return ctx.Status(500).SendString(todoErr.ErrInternalServerError.Error())
	}
	ctx.Response().Header.SetContentType(resp.Header.Get("Content-Type"))

	return ctx.Status(resp.StatusCode).Send(body)

}
