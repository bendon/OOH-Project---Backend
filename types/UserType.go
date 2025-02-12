package types

import "github.com/google/uuid"

type UserPayload struct {
	FirstName  string  `json:"firstName"`
	MiddleName string  `json:"middleName"`
	LastName   string  `json:"lastName"`
	Email      string  `json:"email"`
	UserNumber string  `json:"userNumber"`
	Gender     int     `json:"gender"`
	Phone      int     `json:"phone"`
	Country    *string `json:"country"`
	BirthDate  string  `json:"birthDate"`
	Password   string  `json:"password"`
}

type CreateStaffRequest struct {
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Email      string    `json:"email"`
	Gender     int       `json:"gender"`
	Phone      int       `json:"phone"`
	Country    string    `json:"country"`
	BirthDate  string    `json:"birthDate"`
	RoleId     uuid.UUID `json:"roleId"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SwitchAccountPayload struct {
	AccountId uuid.UUID `json:"accountId"`
	BranchId  uuid.UUID `json:"branchId"`
}

type CreateRoleRequest struct {
	Name string `json:"name"`
}
type UpdateRoleRequest struct {
	RoleId uuid.UUID `json:"roleId"`
	Name   string    `json:"name"`
}

type ChangeEmployeeRoleRequest struct {
	RoleId uuid.UUID `json:"roleId"`
	UserId uuid.UUID `json:"userId"`
}

type ChangeEmployeePermissionRequest struct {
	PermissionIds []uuid.UUID `json:"permissionIds"`
	AccountId     uuid.UUID   `json:"accountId"`
}

type CreateUserBranchRequest struct {
	BranchId uuid.UUID `json:"branchId"`
	UserId   uuid.UUID `json:"userId"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
