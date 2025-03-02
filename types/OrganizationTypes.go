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
	CampaignBrand       *string           `json:"campaignBrand"`
	CampaignDescription string            `json:"campaignDescription" validate:"required"`
	StartDate           *string           `json:"startDate"`
	EndDate             *string           `json:"endDate"`
	Email               *StringArrayJSONB `json:"email" `
	Phone               *Int64ArrayJSONB  `json:"phone" `
	Location            *string           `json:"location" `
	ClientFirstName     *string           `json:"clientFirstName" `
	ClientLastName      *string           `json:"clientLastName" `
	CampaignInsight     *string           `json:"compaignInsight"`
	BillboardId         uuid.UUID         `json:"billboardId" validate:"required"`
	ImageId             *uuid.UUID        `json:"imageId"` // not required
	TargetAudience      *string           `jsonn:"targetAudience"`
	TargetAge           *string           `jsonn:"targetAge"`
	TargetGender        *string           `jsonn:"targetGender"`
	CampainSocials      *StringArrayJSONB `jsonn:"campainSocials"`
	OtherDetails        *JSONB            `jsonn:"otherDetails"`
	CampaignSocials     *StringArrayJSONB `jsonn:"campainSocials"`
	Products            *StringArrayJSONB `jsonn:"products"`
	SiteUrl             *StringArrayJSONB `jsonn:"siteUrl"`
}

type UpdateBillboardCampaignRequest struct {
	CampaignBrand       *string           `json:"campaignBrand"`
	CampaignDescription string            `json:"campaignDescription"`
	StartDate           *string           `json:"startDate"`
	EndDate             *string           `json:"endDate"`
	Email               *StringArrayJSONB `json:"email" `
	Phone               *Int64ArrayJSONB  `json:"phone" `
	Location            *string           `json:"location" `
	ClientFirstName     *string           `json:"clientFirstName" `
	ClientLastName      *string           `json:"clientLastName" `
	CampaignInsight     *string           `json:"compaignInsight"`
	ImageId             *uuid.UUID        `json:"imageId"` // not required
	Others              *JSONB            `jsonn:"others"`
	TargetAudience      *string           `jsonn:"targetAudience"`
	TargetAge           *string           `jsonn:"targetAge"`
	TargetGender        *string           `jsonn:"targetGender"`
	CampainSocials      *StringArrayJSONB `jsonn:"campainSocials"`
	Products            *StringArrayJSONB `jsonn:"products"`
	SiteUrl             *StringArrayJSONB `jsonn:"siteUrl"`
}
