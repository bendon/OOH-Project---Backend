package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillboardCampaignModel struct {
	ID                  uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationId      uuid.UUID       `gorm:"type:char(36);column:organization_id;null" json:"organizationId"`
	CreatedById         uuid.UUID       `gorm:"type:char(36);column:created_by_id; null" json:"createdById"`
	BillboardId         uuid.UUID       `gorm:"type:char(36);column:billboard_id; null" json:"billboardId"`
	Email               *string         `gorm:"type:varchar(255);column:email; null" json:"email"`
	Location            *string         `gorm:"type:varchar(255);column:location;null" json:"location"`
	CampaignBrand       *string         `gorm:"type:varchar(3000);column:campaign_brand;" json:"campaignBrand"`
	CampaignDescription string          `gorm:"type:varchar(3000);column:campaign_description;" json:"campaignDescription"`
	StartDate           *int64          `gorm:"type:bigint;column:start_date;null" json:"startDate"`
	EndDate             *int64          `gorm:"type:bigint;column:end_date;null" json:"endDate"`
	Phone               *int64          `gorm:"type:bigint;column:phone;null" json:"phone"`
	ClientFirstName     *string         `gorm:"type:varchar(255);column:client_first_name;null" json:"clientFirstName"`
	ClientLastName      *string         `gorm:"type:varchar(255);column:client_last_name;null" json:"clientLastName"`
	CampaignInsights    *string         `gorm:"type:varchar(3000);column:campaign_insights;null" json:"campaignInsights"`
	Others              JSONB           `gorm:"type:json;column:others;null" json:"others,omitempty"`
	ImageId             *uuid.UUID      `gorm:"type:char(36);column:image_id; null" json:"image_id"`
	Image               *FileModel      `gorm:"foreignKey:ImageId; references:ID" json:"image,omitempty"`
	Billboard           *BillboardModel `gorm:"foreignKey:BillboardId; references:ID" json:"billboard,omitempty"`
	Active              bool            `gorm:"type:boolean;column:active;default:true;not null" json:"active"`
	ClosedDate          *int64          `gorm:"type:bigint;column:closed_date;null" json:"closedDate"`
	CreatedAt           int64           `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt           int64           `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt           gorm.DeletedAt  `gorm:"index" json:"-"`
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

type JSONB map[string]interface{}

// Marshal JSON before saving to DB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	return json.Marshal(j)
}

// Unmarshal JSON when retrieving from DB
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, &j)
}
