package models

import "github.com/google/uuid"

type OrganizationUserSummary struct {
	OrganizationUserID uuid.UUID  `gorm:"column:organization_user_id" json:"organizationUserId"`
	OrganizationID     uuid.UUID  `gorm:"column:organization_id" json:"organizationId"`
	OrganizationName   string     `gorm:"column:organization_name" json:"organizationName"`
	UserID             uuid.UUID  `gorm:"column:user_id" json:"userId"`
	FirstName          string     `gorm:"column:first_name" json:"firstName"`
	MiddleName         *string    `gorm:"column:middle_name" json:"middleName,omitempty"`
	LastName           string     `gorm:"column:last_name" json:"lastName"`
	Email              string     `gorm:"column:email" json:"email"`
	Phone              int        `gorm:"column:phone" json:"phone"`
	Country            *string    `gorm:"column:country" json:"country,omitempty"`
	Gender             int        `gorm:"column:gender" json:"gender"`
	Verified           bool       `gorm:"column:verified" json:"verified"`
	Active             bool       `gorm:"column:active" json:"active"`
	RoleID             *uuid.UUID `gorm:"column:role_id" json:"roleId,omitempty"`
	RoleName           *string    `gorm:"column:role_name" json:"roleName,omitempty"`
	CreatedAt          int64      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          int64      `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName sets the name of the database view for GORM
func (OrganizationUserSummary) TableName() string {
	return "organization_user_summary"
}
