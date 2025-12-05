package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Category{},
		&Section{},
		&User{},
		&UserAuth{},
		&Course{},
		&CourseRegistration{},
		&AttendanceSession{},
		&AttendanceRecord{},
		&StudentSectionHistory{},
		&Questionnaire{},
		&Question{},
		&QuestionResponse{},
		&QuestionnaireSubmission{},
	)
}
