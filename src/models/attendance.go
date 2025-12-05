package models

import "time"

type AttendanceSession struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CourseID    uint      `json:"courseId" gorm:"not null"`
	Course      Course    `json:"course" gorm:"foreignKey:CourseID"`
	SectionID   *uint     `json:"sectionId"`
	Section     *Section  `json:"section" gorm:"foreignKey:SectionID"`
	SessionDate time.Time `json:"sessionDate" gorm:"type:date;not null"`
	Topic       string    `json:"topic"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (AttendanceSession) TableName() string { return "attendance_sessions" }

type AttendanceRecord struct {
	ID        uint              `json:"id" gorm:"primaryKey"`
	SessionID uint              `json:"sessionId" gorm:"not null;uniqueIndex:idx_session_student"`
	Session   AttendanceSession `json:"session" gorm:"foreignKey:SessionID"`
	StudentID uint              `json:"studentId" gorm:"not null;uniqueIndex:idx_session_student"`
	Student   User              `json:"student" gorm:"foreignKey:StudentID"`
	Status    string            `json:"status" gorm:"type:varchar(20);default:'present'"` // present/absent/late/excused
	Notes     string            `json:"notes"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (AttendanceRecord) TableName() string { return "attendance_records" }
