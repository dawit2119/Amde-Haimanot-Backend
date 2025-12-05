package routes

import (
	"AmdeHaimanot-Backend-Go/src/config"
	"AmdeHaimanot-Backend-Go/src/middleware"
	"AmdeHaimanot-Backend-Go/src/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		page := parseIntDefault(c.Query("page"), 1)
		limit := parseIntDefault(c.Query("limit"), 10)
		offset := (page - 1) * limit

		query := config.DB.Model(&models.User{}).Preload("Category").Preload("Section")
		if cat := c.Query("categoryId"); cat != "" && cat != "all" {
			query = query.Where("category_id = ?", cat)
		}
		if sec := c.Query("sectionId"); sec != "" && sec != "all" {
			query = query.Where("section_id = ?", sec)
		}
		if search := c.Query("search"); search != "" {
			like := "%" + search + "%"
			query = query.Where("full_name ILIKE ? OR phone_no ILIKE ? OR job ILIKE ?", like, like, like)
		}

		var total int64
		query.Count(&total)

		var users []models.User
		query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&users)

		c.JSON(200, gin.H{
			"users":      users,
			"pagination": gin.H{"total": total, "page": page, "totalPages": (total + int64(limit) - 1) / int64(limit), "limit": limit},
		})
	})

	rg.GET("/:id", func(c *gin.Context) {
		var user models.User
		if err := config.DB.Preload("Category").Preload("Section").First(&user, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		c.JSON(200, user)
	})

	rg.GET("/:id/history", func(c *gin.Context) {
		var student models.User
		if err := config.DB.Preload("Category").Preload("Section").First(&student, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Student not found"})
			return
		}

		var registrations []models.CourseRegistration
		config.DB.Preload("Course").Where("student_id = ?", student.ID).Order("created_at DESC").Find(&registrations)

		var attendance []models.AttendanceRecord
		config.DB.Preload("Session.Course").Preload("Session.Section").
			Where("student_id = ?", student.ID).Order("created_at DESC").Find(&attendance)

		var sectionHistory []models.StudentSectionHistory
		config.DB.Preload("FromSection").Preload("ToSection").
			Where("student_id = ?", student.ID).Order("promoted_at DESC").Find(&sectionHistory)

		summary := gin.H{"present": 0, "absent": 0, "late": 0, "excused": 0}
		for _, a := range attendance {
			summary[a.Status] = summary[a.Status].(int) + 1
		}

		c.JSON(200, gin.H{
			"student": student, "registrations": registrations,
			"attendanceRecords": attendance, "attendanceSummary": summary, "sectionHistory": sectionHistory,
		})
	})

	rg.POST("/promote", func(c *gin.Context) {
		var body struct {
			FromSectionID *uint  `json:"fromSectionId"`
			ToSectionID   uint   `json:"toSectionId" binding:"required"`
			StudentIDs    []uint `json:"studentIds"`
			Notes         string `json:"notes"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Target class/section is required"})
			return
		}

		var target models.Section
		if err := config.DB.First(&target, body.ToSectionID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Target section not found"})
			return
		}

		query := config.DB.Model(&models.User{})
		if len(body.StudentIDs) > 0 {
			query = query.Where("id IN ?", body.StudentIDs)
		} else if body.FromSectionID != nil {
			query = query.Where("section_id = ?", *body.FromSectionID)
		} else {
			c.JSON(400, gin.H{"error": "Select students or a source section to promote"})
			return
		}

		var students []models.User
		query.Find(&students)

		for _, student := range students {
			previous := student.SectionID
			config.DB.Model(&student).Update("section_id", body.ToSectionID)
			config.DB.Create(&models.StudentSectionHistory{
				StudentID: student.ID, FromSectionID: previous, ToSectionID: body.ToSectionID, Notes: body.Notes,
			})
		}

		c.JSON(200, gin.H{"message": "Students promoted successfully", "promotedCount": len(students), "toSection": target})
	})

	rg.POST("", func(c *gin.Context) {
		var input middleware.UserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			middleware.RespondValidationError(c, err)
			return
		}
		user := models.User{
			FullName: input.FullName, PhoneNo: input.PhoneNo, CurrentAddress: input.CurrentAddress,
			Job: input.Job, Email: input.Email, CategoryID: input.CategoryID, SectionID: &input.SectionID,
		}
		config.DB.Create(&user)
		config.DB.Preload("Category").Preload("Section").First(&user, user.ID)
		c.JSON(201, user)
	})

	rg.PUT("/:id", func(c *gin.Context) {
		var user models.User
		if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		var input middleware.UserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			middleware.RespondValidationError(c, err)
			return
		}
		config.DB.Model(&user).Updates(models.User{
			FullName: input.FullName, PhoneNo: input.PhoneNo, CurrentAddress: input.CurrentAddress,
			Job: input.Job, Email: input.Email, CategoryID: input.CategoryID, SectionID: &input.SectionID,
		})
		config.DB.Preload("Category").Preload("Section").First(&user, user.ID)
		c.JSON(200, user)
	})

	rg.DELETE("/:id", func(c *gin.Context) {
		var user models.User
		if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		config.DB.Delete(&user)
		c.Status(204)
	})
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
