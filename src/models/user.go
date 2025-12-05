package models

import "time"

type User struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	FullName         string     `json:"fullName" gorm:"size:150;not null"`
	PhoneNo          string     `json:"phoneNo"`
	ProfilePhoto     string     `json:"profilePhoto"`
	CurrentAddress   string     `json:"currentAddress"`
	Job              string     `json:"job"`
	Email            string     `json:"email"`
	BirthDate        *time.Time `json:"birthDate"`
	EmergencyContact string     `json:"emergencyContact"`
	JoinedDate       time.Time  `json:"joinedDate" gorm:"default:CURRENT_DATE"`
	Status           string     `json:"status" gorm:"type:varchar(20);default:'active'"` // enum: active/inactive/pending
	Notes            string     `json:"notes"`

	CategoryID uint     `json:"categoryId" gorm:"not null"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`

	SectionID *uint    `json:"sectionId"`
	Section   *Section `json:"section" gorm:"foreignKey:SectionID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string { return "users" }
