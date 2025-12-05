package models

import "time"

type Category struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	NameAm      string `json:"name_am" gorm:"not null"`
	NameEn      string `json:"name_en" gorm:"not null;unique"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive" gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Category) TableName() string { return "categories" }
