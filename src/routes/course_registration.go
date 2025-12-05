package routes

import (
	"AmdeHaimanot-Backend-Go/src/config"
	"AmdeHaimanot-Backend-Go/src/models"

	"github.com/gin-gonic/gin"
)

func RegisterCourseRegistrationRoutes(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		query := config.DB.Preload("Course").Preload("Student.Category").Preload("Student.Section")
		if courseID := c.Query("courseId"); courseID != "" {
			query = query.Where("course_id = ?", courseID)
		}
		if studentID := c.Query("studentId"); studentID != "" {
			query = query.Where("student_id = ?", studentID)
		}
		if status := c.Query("status"); status != "" {
			query = query.Where("status = ?", status)
		}
		var regs []models.CourseRegistration
		query.Order("created_at DESC").Find(&regs)
		c.JSON(200, regs)
	})

	rg.POST("", func(c *gin.Context) {
		var input struct {
			StudentID uint   `json:"studentId" binding:"required"`
			CourseID  uint   `json:"courseId" binding:"required"`
			Status    string `json:"status"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Validation error"})
			return
		}

		var student models.User
		if err := config.DB.First(&student, input.StudentID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Student not found"})
			return
		}
		var course models.Course
		if err := config.DB.First(&course, input.CourseID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Course not found"})
			return
		}

		status := input.Status
		if status == "" {
			status = "active"
		}

		var reg models.CourseRegistration
		result := config.DB.Where("student_id = ? AND course_id = ?", input.StudentID, input.CourseID).
			FirstOrCreate(&reg, models.CourseRegistration{StudentID: input.StudentID, CourseID: input.CourseID, Status: status})
		if result.Error != nil {
			c.JSON(400, gin.H{"error": "Student is already registered for this course"})
			return
		}
		if reg.Status != status {
			reg.Status = status
			config.DB.Save(&reg)
		}

		config.DB.Preload("Course").Preload("Student.Category").Preload("Student.Section").First(&reg, reg.ID)
		c.JSON(201, reg)
	})

	rg.POST("/bulk", func(c *gin.Context) {
		var input struct {
			CourseID  uint `json:"courseId" binding:"required"`
			SectionID uint `json:"sectionId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Validation error"})
			return
		}

		var course models.Course
		var section models.Section
		var students []models.User
		if err := config.DB.First(&course, input.CourseID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Course not found"})
			return
		}
		if err := config.DB.First(&section, input.SectionID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Section not found"})
			return
		}
		config.DB.Where("section_id = ?", input.SectionID).Find(&students)
		if len(students) == 0 {
			c.JSON(400, gin.H{"error": "No students found in this section"})
			return
		}

		for _, student := range students {
			var reg models.CourseRegistration
			config.DB.Where("student_id = ? AND course_id = ?", student.ID, input.CourseID).
				FirstOrCreate(&reg, models.CourseRegistration{StudentID: student.ID, CourseID: input.CourseID, Status: "active"})
			if reg.Status != "active" {
				reg.Status = "active"
				config.DB.Save(&reg)
			}
		}

		var registered []models.CourseRegistration
		config.DB.Preload("Course").Preload("Student.Category").Preload("Student.Section").
			Where("course_id = ?", input.CourseID).Order("created_at DESC").Find(&registered)

		c.JSON(201, gin.H{
			"message": "Students registered to course", "registeredCount": len(registered),
			"course": course, "section": section, "registrations": registered,
		})
	})

	rg.PUT("/:id", func(c *gin.Context) {
		var reg models.CourseRegistration
		if err := config.DB.First(&reg, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Registration not found"})
			return
		}
		var body struct {
			Status string `json:"status" binding:"omitempty,oneof=active dropped completed"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid registration status"})
			return
		}
		if body.Status != "" {
			reg.Status = body.Status
			config.DB.Save(&reg)
		}
		config.DB.Preload("Course").Preload("Student.Category").Preload("Student.Section").First(&reg, reg.ID)
		c.JSON(200, reg)
	})

	rg.DELETE("/:id", func(c *gin.Context) {
		var reg models.CourseRegistration
		if err := config.DB.First(&reg, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Registration not found"})
			return
		}
		config.DB.Delete(&reg)
		c.Status(204)
	})
}
