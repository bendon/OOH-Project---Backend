package types

import "github.com/google/uuid"

type CreateRoleRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateRoleRequest struct {
	RoleId uuid.UUID `json:"roleId"  validate:"required"`
	Name   string    `json:"name"  validate:"required"`
}

type CreateBillboardCampaignRequest struct {
	CampaignDescription string     `json:"campaignDescription" validate:"required"`
	StartDate           *string    `json:"startDate"`
	EndDate             *string    `json:"endDate"`
	Email               *string    `json:"email" `
	Phone               *int64     `json:"phone" `
	Location            *string    `json:"location" `
	ClientFirstName     *string    `json:"clientFirstName" `
	ClientLastName      *string    `json:"clientLastName" `
	CampaignInsight     *string    `json:"compaignInsight"`
	BillboardId         uuid.UUID  `json:"billboardId" validate:"required"`
	ImageId             *uuid.UUID `json:"imageId"` // not required
}

type UpdateBillboardCampaignRequest struct {
	CampaignDescription string     `json:"campaignDescription"`
	StartDate           *string    `json:"startDate"`
	EndDate             *string    `json:"endDate"`
	Email               *string    `json:"email" `
	Phone               *int64     `json:"phone" `
	Location            *string    `json:"location" `
	ClientFirstName     *string    `json:"clientFirstName" `
	ClientLastName      *string    `json:"clientLastName" `
	CampaignInsight     *string    `json:"compaignInsight"`
	ImageId             *uuid.UUID `json:"imageId"` // not required
}
