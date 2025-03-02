package types

import "github.com/google/uuid"

type UserPayload struct {
	FirstName  string  `json:"firstName"`
	MiddleName string  `json:"middleName"`
	LastName   string  `json:"lastName"`
	Email      string  `json:"email"`
	Gender     int     `json:"gender"`
	Phone      int     `json:"phone"`
	Country    *string `json:"country"`
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
	RoleId     uuid.UUID `json:"roleId"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SwitchAccountPayload struct {
	AccountId uuid.UUID `json:"accountId"`
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

type AuthGoogleVerificationRequest struct {
	Token string `json:"token"`
}

type UserInfo struct {
	Sub           string `json:"sub"`            // Unique Google user ID
	Name          string `json:"name"`           // Full name
	Email         string `json:"email"`          // Email address
	EmailVerified bool   `json:"email_verified"` // Whether the email is verified
	Picture       string `json:"picture"`        // Profile picture URL
	GivenName     string `json:"given_name"`     // First name
	FamilyName    string `json:"family_name"`    // Last name
}

type UpdateStaffPermissionsRequest struct {
	PermissionIds []uuid.UUID `json:"permissionIds"`
}

type UserRegisterRequest struct {
	FirstName  string `json:"firstName" validate:"required"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName" validate:"required"`
	Email      string `json:"email" validate:"required"`
	Gender     *int   `json:"gender"`
	Phone      *int   `json:"phone"`
	Password   string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}
type UpdatePasswordPayload struct {
	Password string `json:"password"`
}
