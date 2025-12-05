package models

import "time"

type Course struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"size:30;not null;unique"`
	Title       string `json:"title" gorm:"size:150;not null"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive" gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Course) TableName() string { return "courses" }

type CourseRegistration struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	StudentID    uint      `json:"studentId" gorm:"not null;uniqueIndex:idx_student_course"`
	Student      User      `json:"student" gorm:"foreignKey:StudentID"`
	CourseID     uint      `json:"courseId" gorm:"not null;uniqueIndex:idx_student_course"`
	Course       Course    `json:"course" gorm:"foreignKey:CourseID"`
	RegisteredAt time.Time `json:"registeredAt" gorm:"default:CURRENT_TIMESTAMP"`
	Status       string    `json:"status" gorm:"type:varchar(20);default:'active'"` // active/dropped/completed
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (CourseRegistration) TableName() string { return "course_registrations" }
