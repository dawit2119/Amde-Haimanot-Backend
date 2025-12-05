package middleware

import (
	"AmdeHaimanot-Backend-Go/src/config"
	"AmdeHaimanot-Backend-Go/src/models"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func parseToken(c *gin.Context, secretEnv string) (*Claims, bool) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(401, gin.H{"error": "No token"})
		return nil, false
	}
	tokenString := strings.TrimPrefix(header, "Bearer ")

	claims := &Claims{}
	secret := os.Getenv(secretEnv)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return nil, false
	}
	return claims, true
}

func VerifyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseToken(c, "JWT_SECRET")
		if !ok {
			c.Abort()
			return
		}
		var user models.UserAuth
		if err := config.DB.First(&user, claims.ID).Error; err != nil || !user.IsActive {
			c.JSON(401, gin.H{"error": "Invalid token or account disabled"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func VerifyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseToken(c, "JWT_SECRET")
		if !ok {
			c.Abort()
			return
		}
		var user models.UserAuth
		if err := config.DB.First(&user, claims.ID).Error; err != nil || user.Role != "admin" {
			c.JSON(403, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
