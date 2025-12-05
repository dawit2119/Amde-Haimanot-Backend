package middleware

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		log.Println(err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Internal server error"})
	}
}
