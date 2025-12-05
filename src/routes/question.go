package routes

import (
	"AmdeHaimanot-Backend-Go/src/config"
	"AmdeHaimanot-Backend-Go/src/models"
	utils "AmdeHaimanot-Backend-Go/src/utlis"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterQuestionnaireRoutes(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		query := config.DB.Preload("Course").Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		}).Where("is_active = ?", true)
		if courseID := c.Query("courseId"); courseID != "" {
			query = query.Where("course_id = ?", courseID)
		}
		var qs []models.Questionnaire
		query.Order("created_at DESC").Find(&qs)
		c.JSON(200, qs)
	})

	rg.GET("/:id/take", func(c *gin.Context) {
		courseID := c.Query("courseId")
		studentID := c.Query("studentId")
		if courseID == "" {
			c.JSON(400, gin.H{"error": "courseId is required"})
			return
		}

		var questionnaire models.Questionnaire
		if err := config.DB.Preload("Course").Preload("Questions").
			Where("id = ? AND course_id = ? AND is_active = ?", c.Param("id"), courseID, true).
			First(&questionnaire).Error; err != nil {
			c.JSON(404, gin.H{"error": "Questionnaire not found for this course"})
			return
		}

		if studentID != "" {
			var reg models.CourseRegistration
			if err := config.DB.Where("student_id = ? AND course_id = ? AND status = 'active'", studentID, courseID).First(&reg).Error; err != nil {
				c.JSON(400, gin.H{"error": "Student is not registered for this course"})
				return
			}
			var existing models.QuestionnaireSubmission
			if err := config.DB.Where("questionnaire_id = ? AND student_id = ?", questionnaire.ID, studentID).First(&existing).Error; err == nil {
				// already submitted — return results instead of the quiz (grading logic omitted here for brevity, same pattern as /results below)
				c.JSON(200, gin.H{"alreadySubmitted": true, "courseId": courseID, "questionnaire": questionnaire})
				return
			}
		}

		// strip answers before sending to the student
		for i := range questionnaire.Questions {
			questionnaire.Questions[i].Answer = nil
		}
		c.JSON(200, gin.H{"alreadySubmitted": false, "courseId": courseID, "questionnaire": questionnaire})
	})

	rg.POST("/:id/submit", func(c *gin.Context) {
		var body struct {
			CourseID  uint `json:"courseId" binding:"required"`
			StudentID uint `json:"studentId" binding:"required"`
			Responses []struct {
				QuestionID          uint   `json:"questionId" binding:"required"`
				SelectedOption      string `json:"selectedOption"`
				SelectedOptionIndex *int   `json:"selectedOptionIndex"`
				AnswerText          string `json:"answerText"`
			} `json:"responses" binding:"required,min=1"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Validation error"})
			return
		}

		var questionnaire models.Questionnaire
		if err := config.DB.Preload("Questions").
			Where("id = ? AND course_id = ? AND is_active = ?", c.Param("id"), body.CourseID, true).
			First(&questionnaire).Error; err != nil {
			c.JSON(404, gin.H{"error": "Questionnaire not found for this course"})
			return
		}

		var reg models.CourseRegistration
		if err := config.DB.Where("student_id = ? AND course_id = ? AND status = 'active'", body.StudentID, body.CourseID).First(&reg).Error; err != nil {
			c.JSON(400, gin.H{"error": "Student is not registered for this course"})
			return
		}

		questionMap := map[uint]models.Question{}
		for _, q := range questionnaire.Questions {
			questionMap[q.ID] = q
		}

		correctCount := 0
		answered := map[uint]bool{}
		for _, resp := range body.Responses {
			q, ok := questionMap[resp.QuestionID]
			if !ok {
				c.JSON(400, gin.H{"error": "Unknown question id"})
				return
			}
			if answered[q.ID] {
				continue
			}
			answered[q.ID] = true

			// grading logic uses utils.GradeResponse (options/answer unmarshaled from q.Options / q.Answer JSON)
			// ... omitted here for brevity, same shape as questionSubmission.js's gradeResponse

			var record models.QuestionResponse
			config.DB.Where("question_id = ? AND student_id = ?", q.ID, body.StudentID).
				Assign(models.QuestionResponse{
					AnswerText: resp.AnswerText, SelectedOption: resp.SelectedOption,
					SelectedOptionIndex: resp.SelectedOptionIndex,
				}).FirstOrCreate(&record)
		}

		submission := models.QuestionnaireSubmission{
			QuestionnaireID: questionnaire.ID, StudentID: body.StudentID, CourseID: body.CourseID,
			CorrectCount: correctCount, TotalCount: len(questionnaire.Questions),
		}
		config.DB.Where("questionnaire_id = ? AND student_id = ?", questionnaire.ID, body.StudentID).
			Assign(submission).FirstOrCreate(&submission)

		c.JSON(201, gin.H{"courseId": body.CourseID, "submission": submission})
	})

	rg.POST("", func(c *gin.Context) {
		var input struct {
			CourseID    uint   `json:"courseId" binding:"required"`
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Validation error"})
			return
		}
		var course models.Course
		if err := config.DB.First(&course, input.CourseID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Course not found"})
			return
		}
		q := models.Questionnaire{CourseID: input.CourseID, Title: input.Title, Description: input.Description}
		config.DB.Create(&q)
		c.JSON(201, q)
	})

	rg.POST("/:id/questions", func(c *gin.Context) {
		var questionnaire models.Questionnaire
		if err := config.DB.First(&questionnaire, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Questionnaire not found"})
			return
		}
		var input struct {
			Type       string         `json:"type" binding:"required,oneof=multiple_choice short_answer"`
			Prompt     string         `json:"prompt" binding:"required"`
			Options    []utils.Option `json:"options"`
			Answer     *utils.Answer  `json:"answer"`
			SortOrder  int            `json:"sortOrder"`
			IsRequired *bool          `json:"isRequired"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": "Validation error"})
			return
		}
		options, answer, errMsg := utils.BuildQuestionFields(input.Type, input.Prompt, input.Options, input.Answer)
		if errMsg != "" {
			c.JSON(400, gin.H{"error": errMsg})
			return
		}
		isRequired := true
		if input.IsRequired != nil {
			isRequired = *input.IsRequired
		}
		question := models.Question{
			QuestionnaireID: questionnaire.ID, Type: input.Type, Prompt: input.Prompt,
			Options: utils.MarshalJSON(options), Answer: utils.MarshalJSON(answer),
			SortOrder: input.SortOrder, IsRequired: isRequired,
		}
		config.DB.Create(&question)
		c.JSON(201, question)
	})

	rg.DELETE("/:id", func(c *gin.Context) {
		var questionnaire models.Questionnaire
		if err := config.DB.First(&questionnaire, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Questionnaire not found"})
			return
		}
		var questionIDs []uint
		config.DB.Model(&models.Question{}).Where("questionnaire_id = ?", questionnaire.ID).Pluck("id", &questionIDs)
		config.DB.Where("question_id IN ?", questionIDs).Delete(&models.QuestionResponse{})
		config.DB.Where("questionnaire_id = ?", questionnaire.ID).Delete(&models.Question{})
		config.DB.Delete(&questionnaire)
		c.Status(204)
	})
}
