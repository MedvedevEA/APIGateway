package dto

import "github.com/google/uuid"

type AddUser struct {
	UserId *uuid.UUID `json:"user_id"`
	Name   string     `json:"name"`
}
