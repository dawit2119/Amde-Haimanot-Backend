package middleware

import "github.com/gin-gonic/gin"

type UserInput struct {
	FullName       string `json:"fullName" binding:"required,min=2,max=150"`
	PhoneNo        string `json:"phoneNo"`
	CurrentAddress string `json:"currentAddress"`
	Job            string `json:"job"`
	Email          string `json:"email"`
	BirthDate      string `json:"birthDate"`
	CategoryID     uint   `json:"categoryId" binding:"required,min=1"`
	SectionID      uint   `json:"sectionId" binding:"required,min=1"`
}

type CategoryInput struct {
	NameAm string `json:"name_am" binding:"required,min=2,max=100"`
	NameEn string `json:"name_en" binding:"required,min=2,max=100"`
}

// shared helper for reporting bind/validation errors consistently
func RespondValidationError(c *gin.Context, err error) {
	c.JSON(400, gin.H{"error": "Validation error", "details": err.Error()})
}
