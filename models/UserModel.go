package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	FirstName  string    `gorm:"column:first_name;not null" json:"firstName"`
	MiddleName string    `gorm:"column:middle_name" json:"middleName"`
	LastName   string    `gorm:"column:last_name" json:"lastName"`
	Email      string    `gorm:"column:email;unique;not null" json:"email"`
	Phone      int       `gorm:"column:phone;unique;not null" json:"phone"`
	Country    string    `gorm:"column:country;null" json:"country"`
	Gender     int       `gorm:"type:int;column:gender;not null" json:"gender"`
	Verified   bool      `gorm:"type:boolean;default:false;column:verified" json:"verified"`
	IsChange   *bool     `gorm:"type:boolean;column:is_change;null" json:"isChange,omitempty"`
	Active     bool      `gorm:"type:boolean;default:false;column:active" json:"active"`
	Password   string    `gorm:"column:password" json:"-"`
	CreatedAt  int64     `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  int64     `gorm:"column:updated_at" json:"updatedAt"`
}

func (UserModel) TableName() string {
	return "users"
}

func (user *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	return
}

func (user *UserModel) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now().Unix()
	return
}
