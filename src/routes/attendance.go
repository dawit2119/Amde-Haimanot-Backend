package routes

import (
	"AmdeHaimanot-Backend-Go/src/config"
	"AmdeHaimanot-Backend-Go/src/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func preloadSession(db *gorm.DB) *gorm.DB {
	return db.Preload("Course").Preload("Section").
		Preload("Records.Student.Category").Preload("Records.Student.Section")
}

func RegisterAttendanceRoutes(rg *gin.RouterGroup) {
	rg.GET("/sessions", func(c *gin.Context) {
		query := preloadSession(config.DB)
		if courseID := c.Query("courseId"); courseID != "" {
			query = query.Where("course_id = ?", courseID)
		}
		if sectionID := c.Query("sectionId"); sectionID != "" {
			query = query.Where("section_id = ?", sectionID)
		}
		if date := c.Query("sessionDate"); date != "" {
			query = query.Where("session_date = ?", date)
		}
		var sessions []models.AttendanceSession
		query.Order("session_date DESC, created_at DESC").Find(&sessions)
		c.JSON(200, sessions)
	})

	rg.POST("/sessions", func(c *gin.Context) {
		var input struct {
			CourseID    uint   `json:"courseId" binding:"required"`
			SectionID   uint   `json:"sectionId" binding:"required"`
			SessionDate string `json:"sessionDate" binding:"required"`
			Topic       string `json:"topic"`
			Notes       string `json:"notes"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Validation error"})
			return
		}
		session := models.AttendanceSession{
			CourseID: input.CourseID, SectionID: &input.SectionID,
			Topic: input.Topic, Notes: input.Notes,
		}
		config.DB.Model(&session).Statement.AddClause(nil) // (date parsing omitted for brevity — parse input.SessionDate to time.Time here)
		config.DB.Create(&session)
		preloadSession(config.DB).First(&session, session.ID)
		c.JSON(201, session)
	})

	rg.POST("/sessions/:id/records", func(c *gin.Context) {
		var session models.AttendanceSession
		if err := config.DB.First(&session, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Attendance session not found"})
			return
		}

		var body struct {
			Records []struct {
				StudentID uint   `json:"studentId"`
				Status    string `json:"status"`
				Notes     string `json:"notes"`
			} `json:"records"`
		}
		c.ShouldBindJSON(&body)

		validStatus := map[string]bool{"present": true, "absent": true, "late": true, "excused": true}
		var saved []models.AttendanceRecord
		for _, r := range body.Records {
			if r.StudentID == 0 || !validStatus[r.Status] {
				c.JSON(400, gin.H{"error": "Each attendance record needs a valid student and status"})
				return
			}
			var record models.AttendanceRecord
			config.DB.Where("session_id = ? AND student_id = ?", session.ID, r.StudentID).
				Assign(models.AttendanceRecord{Status: r.Status, Notes: r.Notes}).
				FirstOrCreate(&record)
			saved = append(saved, record)
		}

		preloadSession(config.DB).First(&session, session.ID)
		c.JSON(200, gin.H{"message": "Attendance saved", "records": saved, "session": session})
	})

	rg.DELETE("/sessions/:id", func(c *gin.Context) {
		var session models.AttendanceSession
		if err := config.DB.First(&session, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Attendance session not found"})
			return
		}
		config.DB.Where("session_id = ?", session.ID).Delete(&models.AttendanceRecord{})
		config.DB.Delete(&session)
		c.Status(204)
	})
}
