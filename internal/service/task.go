package service

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Service) AddTask(ctx *fiber.Ctx) error {
	return s.todolist.Redirect(ctx, "tasks")
}

func (s *Service) GetTask(ctx *fiber.Ctx) error {
	return s.todolist.Redirect(ctx, "tasks")
}
func (s *Service) GetTasks(ctx *fiber.Ctx) error {
	return s.todolist.Redirect(ctx, "tasks")

}
func (s *Service) UpdateTask(ctx *fiber.Ctx) error {
	return s.todolist.Redirect(ctx, "tasks")
}
func (s *Service) RemoveTask(ctx *fiber.Ctx) error {
	return s.todolist.Redirect(ctx, "tasks")
}
