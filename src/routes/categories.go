package routes

import (
	"AmdeHaimanot-Backend-Go/src/config"
	"AmdeHaimanot-Backend-Go/src/middleware"
	"AmdeHaimanot-Backend-Go/src/models"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(rg *gin.RouterGroup) {
	rg.GET("", func(c *gin.Context) {
		var categories []models.Category
		config.DB.Where("is_active = ?", true).Order("name_en ASC").Find(&categories)
		c.JSON(200, categories)
	})

	rg.GET("/with-count", func(c *gin.Context) {
		type result struct {
			models.Category
			MemberCount int64 `json:"memberCount"`
		}
		var categories []models.Category
		config.DB.Where("is_active = ?", true).Order("name_en ASC").Find(&categories)

		var out []result
		for _, cat := range categories {
			var count int64
			config.DB.Model(&models.User{}).Where("category_id = ?", cat.ID).Count(&count)
			out = append(out, result{Category: cat, MemberCount: count})
		}
		c.JSON(200, out)
	})

	rg.POST("", func(c *gin.Context) {
		var input middleware.CategoryInput
		if err := c.ShouldBindJSON(&input); err != nil {
			middleware.RespondValidationError(c, err)
			return
		}
		category := models.Category{NameAm: input.NameAm, NameEn: input.NameEn}
		if err := config.DB.Create(&category).Error; err != nil {
			c.JSON(400, gin.H{"error": "Category with this English name already exists"})
			return
		}
		c.JSON(201, category)
	})

	rg.PUT("/:id", func(c *gin.Context) {
		var category models.Category
		if err := config.DB.First(&category, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Category not found"})
			return
		}
		var input middleware.CategoryInput
		if err := c.ShouldBindJSON(&input); err != nil {
			middleware.RespondValidationError(c, err)
			return
		}
		category.NameAm, category.NameEn = input.NameAm, input.NameEn
		if err := config.DB.Save(&category).Error; err != nil {
			c.JSON(400, gin.H{"error": "Category with this English name already exists"})
			return
		}
		c.JSON(200, category)
	})

	rg.DELETE("/:id", func(c *gin.Context) {
		var category models.Category
		if err := config.DB.First(&category, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "Category not found"})
			return
		}
		var userCount int64
		config.DB.Model(&models.User{}).Where("category_id = ?", category.ID).Count(&userCount)
		if userCount > 0 {
			c.JSON(400, gin.H{"error": "Cannot delete category with associated members. Please reassign or delete members first."})
			return
		}
		config.DB.Delete(&category)
		c.Status(204)
	})
}
