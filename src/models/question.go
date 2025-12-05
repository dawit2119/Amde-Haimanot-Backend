package models

import (
	"time"

	"gorm.io/datatypes"
)

type Questionnaire struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CourseID    uint       `json:"courseId" gorm:"not null"`
	Course      Course     `json:"course" gorm:"foreignKey:CourseID"`
	Title       string     `json:"title" gorm:"size:200;not null"`
	Description string     `json:"description"`
	IsActive    bool       `json:"isActive" gorm:"default:true"`
	Questions   []Question `json:"questions" gorm:"foreignKey:QuestionnaireID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Questionnaire) TableName() string { return "questionnaires" }

type Question struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	QuestionnaireID uint           `json:"questionnaireId" gorm:"not null"`
	Questionnaire   Questionnaire  `json:"questionnaire" gorm:"foreignKey:QuestionnaireID"`
	Type            string         `json:"type" gorm:"type:varchar(30);not null"` // multiple_choice/short_answer
	Prompt          string         `json:"prompt" gorm:"type:text;not null"`
	Options         datatypes.JSON `json:"options"` // [{ text: string }, ...]
	Answer          datatypes.JSON `json:"answer"`  // { text, optionIndex } or { text }
	SortOrder       int            `json:"sortOrder" gorm:"default:0"`
	IsRequired      bool           `json:"isRequired" gorm:"default:true"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Question) TableName() string { return "questions" }

type QuestionResponse struct {
	ID                  uint     `json:"id" gorm:"primaryKey"`
	QuestionID          uint     `json:"questionId" gorm:"not null;uniqueIndex:idx_question_student"`
	Question            Question `json:"question" gorm:"foreignKey:QuestionID"`
	StudentID           uint     `json:"studentId" gorm:"not null;uniqueIndex:idx_question_student"`
	Student             User     `json:"student" gorm:"foreignKey:StudentID"`
	AnswerText          string   `json:"answerText"`
	SelectedOption      string   `json:"selectedOption" gorm:"size:500"`
	SelectedOptionIndex *int     `json:"selectedOptionIndex"`
	IsCorrect           *bool    `json:"isCorrect"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (QuestionResponse) TableName() string { return "question_responses" }

type QuestionnaireSubmission struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	QuestionnaireID uint          `json:"questionnaireId" gorm:"not null;uniqueIndex:idx_questionnaire_student"`
	Questionnaire   Questionnaire `json:"questionnaire" gorm:"foreignKey:QuestionnaireID"`
	StudentID       uint          `json:"studentId" gorm:"not null;uniqueIndex:idx_questionnaire_student"`
	Student         User          `json:"student" gorm:"foreignKey:StudentID"`
	CourseID        uint          `json:"courseId" gorm:"not null"`
	Course          Course        `json:"course" gorm:"foreignKey:CourseID"`
	CorrectCount    int           `json:"correctCount" gorm:"default:0"`
	TotalCount      int           `json:"totalCount" gorm:"default:0"`
	SubmittedAt     time.Time     `json:"submittedAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (QuestionnaireSubmission) TableName() string { return "questionnaire_submissions" }
