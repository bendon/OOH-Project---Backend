package models

import (
	"github.com/google/uuid"

	"bbscout/models"
	"bbscout/types"
)

type BillboardSummaryView struct {
	OrganizationId  uuid.UUID                      `gorm:"type:char(36);column:organization_id;primaryKey" json:"organizationId"`
	BillboardId     uuid.UUID                      `gorm:"type:char(36);column:billboard_id;primaryKey" json:"billboardId"`
	CreatedById     uuid.UUID                      `gorm:"type:char(36);column:created_by_id" json:"createdById"`
	BoardCode       string                         `gorm:"type:varchar(1000);column:board_code; null" json:"boardCode"`
	Location        string                         `gorm:"type:varchar(255);column:location;not null" json:"location"`
	Latitude        float64                        `gorm:"type:double;column:latitude;not null" json:"latitude"`
	Longitude       float64                        `gorm:"type:double;column:longitude;not null" json:"longitude"`
	Width           float64                        `gorm:"type:double;column:width;not null" json:"width"`
	Height          float64                        `gorm:"type:double;column:height;not null" json:"height"`
	Unit            string                         `gorm:"type:enum('centimeters','meters','feet','inches');column:unit;not null" json:"unit"`
	Type            string                         `gorm:"type:enum('digital','static','LED','traditional');column:type;not null" json:"type"`
	Price           *float64                       `gorm:"type:decimal(10,2);column:price;null;default:0" json:"price"`
	BillboardActive bool                           `gorm:"type:boolean;column:billboard_active;not null" json:"billboardActive"`
	CampaignId      *uuid.UUID                     `gorm:"type:char(36);column:campaign_id" json:"campaignId"`
	CampaignActive  bool                           `gorm:"type:boolean;column:campaign_active;not null;default:false" json:"campaignActive"`
	ImageId         *uuid.UUID                     `gorm:"type:char(36);column:image_id; null" json:"image_id"`
	Image           *models.FileModel              `gorm:"foreignKey:ImageId; references:ID" json:"image,omitempty"`
	Staff           *models.UserModel              `gorm:"foreignKey:CreatedById; references:ID" json:"staff,omitempty"`
	Campaign        *models.BillboardCampaignModel `gorm:"foreignKey:CampaignId; references:ID" json:"campaign,omitempty"`
	CreatedAt       int64                          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       int64                          `gorm:"column:updated_at" json:"updatedAt"`
	Occupied        bool                           `gorm:"type:boolean;column:occupied;not null" json:"occupied"`
	Owner           *string                        `gorm:"type:varchar(255);column:owner;null" json:"owner"`
	OwnerContact    *types.Int64ArrayJSONB         `gorm:"type:json;column:owner_contact;null" json:"ownerContact"`
	OwnerEmail      *types.StringArrayJSONB        `gorm:"type:json;column:owner_email;null" json:"ownerEmail"`
	CloseupImageId  *uuid.UUID                     `gorm:"type:char(36);column:closeup_image_id;null" json:"closeupImageId"`
	CloseupImage    *models.FileModel              `gorm:"foreignKey:CloseupImageId;references:ID" json:"closeupImage,omitempty"`
	Material        *string                        `gorm:"type:varchar(255);column:material;null" json:"material"`
	Visibility      *string                        `gorm:"type:varchar(255);column:visibility;null" json:"visibility"`
	Illumination    *string                        `gorm:"type:varchar(255);column:illumination;null" json:"illumination"`
	Angle           *string                        `gorm:"type:varchar(255);column:angle;null" json:"angle"`
	Structure       *string                        `gorm:"type:varchar(255);column:structure;null" json:"structure"`
	City            *string                        `gorm:"type:varchar(255);column:city;null" json:"city"`
	ObjectType      *string                        `gorm:"type:varchar(255);column:object_type;null" json:"objectType"`
	ParentBoardCode *string                        `gorm:"type:varchar(1000);column:parent_board_code;null" json:"parentBoardCode"`
}

// TableName overrides the default table name for GORM
func (BillboardSummaryView) TableName() string {
	return "billboard_summary"
}
