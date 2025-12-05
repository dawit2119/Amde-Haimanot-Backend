package routes

import (
	"AmdeHaimanot-Backend-Go/src/config"
	"AmdeHaimanot-Backend-Go/src/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func getField(row map[string]string, keys []string) string {
	for _, k := range keys {
		if v, ok := row[k]; ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func RegisterImportRoutes(rg *gin.RouterGroup) {
	rg.POST("", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "No file uploaded"})
			return
		}
		opened, _ := file.Open()
		defer opened.Close()

		f, err := excelize.OpenReader(opened)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid Excel file"})
			return
		}
		sheet := f.GetSheetList()[0]
		rows, _ := f.GetRows(sheet)
		if len(rows) < 2 {
			c.JSON(400, gin.H{"error": "Excel file is empty"})
			return
		}

		headers := rows[0]
		var categories []models.Category
		config.DB.Find(&categories)

		results := gin.H{"total": len(rows) - 1, "succeeded": 0, "failed": 0}
		var errorsList []gin.H
		succeeded := 0

		for _, r := range rows[1:] {
			row := map[string]string{}
			for i, h := range headers {
				if i < len(r) {
					row[h] = r[i]
				}
			}

			categoryName := getField(row, []string{"የትውልድ መንደር", "category", "Category", "ምድብ", "cat"})
			fullName := getField(row, []string{"fullName", "Full Name", "ሙሉ ስም"})

			if categoryName == "" || fullName == "" {
				errorsList = append(errorsList, gin.H{"row": fullName, "error": "Missing category or full name"})
				continue
			}

			var matchedCategory *models.Category
			for i := range categories {
				if strings.EqualFold(categories[i].NameEn, categoryName) || strings.EqualFold(categories[i].NameAm, categoryName) {
					matchedCategory = &categories[i]
					break
				}
			}
			if matchedCategory == nil {
				errorsList = append(errorsList, gin.H{"row": fullName, "error": "Category not found: " + categoryName})
				continue
			}

			user := models.User{
				FullName: fullName, CategoryID: matchedCategory.ID,
				PhoneNo: getField(row, []string{"phoneNo", "Phone No", "ስልክ ቁጥር"}),
				Email:   getField(row, []string{"email", "Email"}),
			}
			if err := config.DB.Create(&user).Error; err != nil {
				errorsList = append(errorsList, gin.H{"row": fullName, "error": err.Error()})
				continue
			}
			succeeded++
		}

		results["succeeded"] = succeeded
		results["failed"] = len(errorsList)
		results["errors"] = errorsList
		c.JSON(200, results)
	})
}
