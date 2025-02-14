package models

import "github.com/google/uuid"

type OrganizationUserAnalytics struct {
	OrganizationID   uuid.UUID `gorm:"column:organization_id" json:"organizationId"`
	OrganizationName string    `gorm:"column:organization_name" json:"organizationName"`
	NoOfUsers        int       `gorm:"column:no_of_users" json:"noOfUsers"`
	MaleCount        int       `gorm:"column:male_count" json:"maleCount"`
	FemaleCount      int       `gorm:"column:female_count" json:"femaleCount"`
	TransgenderCount int       `gorm:"column:transgender_count" json:"transgenderCount"`
	NoOfRoles        int       `gorm:"column:no_of_roles" json:"noOfRoles"`
	JoinedThisMonth  int       `gorm:"column:joined_this_month" json:"joinedThisMonth"`
	VerifiedUsers    int       `gorm:"column:verified_users" json:"verifiedUsers"`
	ActiveUsers      int       `gorm:"column:active_users" json:"activeUsers"`
	LastUserJoinedAt int64     `gorm:"column:last_user_joined_at" json:"lastUserJoinedAt"`
}

// TableName sets the table name for GORM
func (OrganizationUserAnalytics) TableName() string {
	return "organization_user_analytics"
}
