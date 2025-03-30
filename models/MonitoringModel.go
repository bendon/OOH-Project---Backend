package models

import (
	"bbscout/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MonitoringModel struct {
	ID                   uuid.UUID               `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationId       uuid.UUID               `gorm:"type:char(36);column:organization_id;null" json:"organizationId"`
	MonitoredById        uuid.UUID               `gorm:"type:char(36);column:monitored_by_id; null" json:"monitoredById"`
	BillboardId          *uuid.UUID              `gorm:"type:char(36);column:billboard_id; null" json:"billboardId"`
	Date                 *string                 `gorm:"type:varchar(1000);column:board_code; null" json:"boardCode"`
	County               *string                 `gorm:"type:varchar(255);column:owner; null" json:"Owner"`
	Street               *string                 `gorm:"type:varchar(255);column:street; null" json:"street"`
	Location             *string                 `gorm:"type:varchar(255);column:location; null" json:"location"`
	Building             *string                 `gorm:"type:varchar(255);column:building; null" json:"building"`
	OwnerContacts        *types.Int64ArrayJSONB  `gorm:"type:json;column:owner_contact;null" json:"ownerContacts"`
	OwnerEmails          *types.StringArrayJSONB `gorm:"type:json;column:owner_email;null" json:"ownerEmails"`
	Brand                *string                 `gorm:"type:varchar(255);column:brand; null" json:"brand"`
	Campain              *string                 `gorm:"type:varchar(255);column:campaign; null" json:"campain"`
	Width                *float64                `gorm:"type:double;column:width; null;default:0" json:"width"`
	Height               *float64                `gorm:"type:double;column:height; null;default:0" json:"height"`
	Unit                 *string                 `gorm:"type:enum('centimeters','meters','feet','inches');column:unit; null;default:meters" json:"unit"`
	Material             *string                 `gorm:"type:varchar(255);column:material; null" json:"material"`
	Angle                *string                 `gorm:"type:varchar(255);column:angel; null" json:"angle"`
	Visibility           *string                 `gorm:"type:varchar(255);column:visibility; null" json:"visibility"`
	Illumination         *string                 `gorm:"type:varchar(255);column:illumination;null" json:"illumination"`
	Environment          *string                 `gorm:"type:varchar(255);column:environment; null" json:"environment"`
	ConditionOfMaterial  *string                 `gorm:"type:varchar(255);column:condition_of_material;null" json:"conditionOfMaterial"`
	ConditionOfStructure *string                 `gorm:"type:varchar(255);column:condition_of_structure;null" json:"conditionOfStructure"`
	Comments             *string                 `gorm:"type:varchar(255);column:comments; null" json:"comments"`
	Latitude             *float64                `gorm:"type:double;column:latitude; null; default:0" json:"latitude"`
	Longitude            *float64                `gorm:"type:double;column:longitude; null; default:0" json:"longitude"`
	Accuracy             *float64                `gorm:"type:double;column:accuracy;null; default:0" json:"accuracy"`
	LongShotImageId      *uuid.UUID              `gorm:"type:char(36);column:image_id; null" json:"imageId"`
	CloseUpImageId       *uuid.UUID              `gorm:"type:char(36);column:closeup_image_id; null" json:"closeUpImageId"`
	Structure            *string                 `gorm:"type:varchar(255);column:structure; null" json:"structure"`
	Type                 *string                 `gorm:"type:varchar(255);column:type; null" json:"type"` // "Static Billboard", "Digital Billboard", "Banner Ads", "Wallscapes", "Mobile Billboards","Lamp Posts","Interactive Billboards"
	CreatedAt            int64                   `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt            int64                   `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt            gorm.DeletedAt          `gorm:"index" json:"-"`
	User                 *UserModel              `gorm:"foreignKey:MonitoredById; references:ID" json:"user,omitempty"`
	LongShortImage       *FileModel              `gorm:"foreignKey:LongShotImageId; references:ID" json:"longShotImage,omitempty"`
	CloseUpImage         *FileModel              `gorm:"foreignKey:CloseUpImageId; references:ID" json:"closeUpImage,omitempty"`
	Billboard            *BillboardModel         `gorm:"foreignKey:BillboardById; references:ID" json:"billboard,omitempty"`
}

func (MonitoringModel) TableName() string {
	return "monitoring_strutcures"
}

func (b *MonitoringModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
	return
}

func (b *MonitoringModel) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now().Unix()
	return

}
