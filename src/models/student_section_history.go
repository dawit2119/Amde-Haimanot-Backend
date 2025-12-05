package models

import "time"

type StudentSectionHistory struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	StudentID     uint      `json:"studentId" gorm:"not null"`
	Student       User      `json:"student" gorm:"foreignKey:StudentID"`
	FromSectionID *uint     `json:"fromSectionId"`
	FromSection   *Section  `json:"fromSection" gorm:"foreignKey:FromSectionID"`
	ToSectionID   uint      `json:"toSectionId" gorm:"not null"`
	ToSection     Section   `json:"toSection" gorm:"foreignKey:ToSectionID"`
	PromotedAt    time.Time `json:"promotedAt" gorm:"default:CURRENT_TIMESTAMP"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (StudentSectionHistory) TableName() string { return "student_section_history" }
