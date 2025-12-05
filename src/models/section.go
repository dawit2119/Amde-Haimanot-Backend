package models

import "time"

type Section struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"size:20;not null;unique"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive" gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Section) TableName() string { return "sections" }
