package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"bbscout/types"
)

type BillboardCampaignModel struct {
	ID                  uuid.UUID               `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationId      uuid.UUID               `gorm:"type:char(36);column:organization_id;null" json:"organizationId"`
	CreatedById         uuid.UUID               `gorm:"type:char(36);column:created_by_id; null" json:"createdById"`
	BillboardId         uuid.UUID               `gorm:"type:char(36);column:billboard_id; null" json:"billboardId"`
	Email               *types.StringArrayJSONB `gorm:"type:json;column:email; null" json:"email"`
	Location            *string                 `gorm:"type:varchar(255);column:location;null" json:"location"`
	CampaignBrand       *string                 `gorm:"type:varchar(3000);column:campaign_brand;" json:"campaignBrand"`
	CampaignDescription string                  `gorm:"type:varchar(3000);column:campaign_description;" json:"campaignDescription"`
	StartDate           *int64                  `gorm:"type:bigint;column:start_date;null" json:"startDate"`
	EndDate             *int64                  `gorm:"type:bigint;column:end_date;null" json:"endDate"`
	Phone               *types.Int64ArrayJSONB  `gorm:"type:json;column:phone;null" json:"phone"`
	ClientFirstName     *string                 `gorm:"type:varchar(255);column:client_first_name;null" json:"clientFirstName"`
	ClientLastName      *string                 `gorm:"type:varchar(255);column:client_last_name;null" json:"clientLastName"`
	CampaignInsights    *string                 `gorm:"type:varchar(3000);column:campaign_insights;null" json:"campaignInsights"`
	Others              *types.JSONB            `gorm:"type:json;column:others;null" json:"others,omitempty"`
	ImageId             *uuid.UUID              `gorm:"type:char(36);column:image_id; null" json:"image_id"`
	Image               *FileModel              `gorm:"foreignKey:ImageId; references:ID" json:"image,omitempty"`
	Billboard           *BillboardModel         `gorm:"foreignKey:BillboardId; references:ID" json:"billboard,omitempty"`
	Active              bool                    `gorm:"type:boolean;column:active;default:true;not null" json:"active"`
	ClosedDate          *int64                  `gorm:"type:bigint;column:closed_date;null" json:"closedDate"`
	SiteUrl             *types.StringArrayJSONB `gorm:"type:json;column:site_url;" json:"siteUrl"`
	TargetAudience      *string                 `gorm:"type:varchar(1000);column:target_audience;" json:"targetAudience"`
	TargetAge           *string                 `gorm:"type:varchar(255);column:target_age;" json:"targetAge"`
	TargetGender        *string                 `gorm:"type:varchar(255);column:target_gender;" json:"targetGender"`
	CampaignSocials     *types.StringArrayJSONB `gorm:"type:json;column:campaign_socials;" json:"campaignSocials"`
	Products            *types.StringArrayJSONB `gorm:"type:json;column:products;" json:"products"`
	ImageUrl            *string                 `gorm:"type:varchar(255);column:image_url;" json:"imageUrl"`
	CreatedAt           int64                   `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt           int64                   `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt           gorm.DeletedAt          `gorm:"index" json:"-"`
}

func (BillboardCampaignModel) TableName() string {
	return "billboard_campaign"
}

func (b *BillboardCampaignModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
	return
}

func (b *BillboardCampaignModel) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now().Unix()
	return

}
