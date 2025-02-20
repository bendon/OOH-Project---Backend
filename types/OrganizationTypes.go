package types

import "github.com/google/uuid"

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateRoleRequest struct {
	RoleId uuid.UUID `json:"roleId"  validate:"required"`
	Name   string    `json:"name"  validate:"required"`
}
