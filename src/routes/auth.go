package routes

import (
	"AmdeHaimanot-Backend-Go/src/config"
	"AmdeHaimanot-Backend-Go/src/middleware"
	"AmdeHaimanot-Backend-Go/src/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullName" binding:"required"`
	Role     string `json:"role" binding:"omitempty,oneof=admin manager viewer"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func generateTokens(user models.UserAuth) (string, string, error) {
	claims := middleware.Claims{ID: user.ID, Username: user.Username, Role: user.Role}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": claims.ID, "username": claims.Username, "role": claims.Role,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	accessStr, err := access.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", err
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": claims.ID, "username": claims.Username, "role": claims.Role,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	refreshStr, err := refresh.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	return accessStr, refreshStr, err
}

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"errors": err.Error()})
			return
		}

		var existing models.UserAuth
		if err := config.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existing).Error; err == nil {
			c.JSON(400, gin.H{"error": "Username or email already exists"})
			return
		}

		var count int64
		config.DB.Model(&models.UserAuth{}).Count(&count)
		newUserNumber := count + 1

		assignedRole := "viewer"
		if newUserNumber == 1 || newUserNumber == 3 || newUserNumber == 4 {
			assignedRole = "admin"
		} else if input.Role != "" {
			assignedRole = input.Role
		}

		user := models.UserAuth{
			Username: input.Username, Email: input.Email,
			Password: input.Password, FullName: input.FullName, Role: assignedRole,
		}
		if err := user.HashPassword(); err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		access, refresh, err := generateTokens(user)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(201, gin.H{
			"message": "User registered successfully",
			"access":  access, "refresh": refresh,
			"user": gin.H{"id": user.ID, "username": user.Username, "email": user.Email, "fullName": user.FullName, "role": user.Role},
		})
	})

	rg.POST("/login", func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"errors": err.Error()})
			return
		}

		var user models.UserAuth
		if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil || !user.ValidatePassword(input.Password) {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}
		if !user.IsActive {
			c.JSON(401, gin.H{"error": "Account disabled"})
			return
		}

		access, refresh, _ := generateTokens(user)
		c.JSON(200, gin.H{
			"access": access, "refresh": refresh,
			"user": gin.H{"id": user.ID, "username": user.Username, "email": user.Email, "fullName": user.FullName, "role": user.Role},
		})
	})

	rg.POST("/refresh", func(c *gin.Context) {
		var body struct {
			Refresh string `json:"refresh" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(401, gin.H{"error": "Refresh token required"})
			return
		}

		claims := &middleware.Claims{}
		token, err := jwt.ParseWithClaims(body.Refresh, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid refresh token"})
			return
		}

		var user models.UserAuth
		if err := config.DB.First(&user, claims.ID).Error; err != nil || !user.IsActive {
			c.JSON(401, gin.H{"error": "Invalid refresh token"})
			return
		}

		access, refresh, _ := generateTokens(user)
		c.JSON(200, gin.H{"access": access, "refresh": refresh})
	})

	rg.GET("/me", middleware.VerifyAuth(), func(c *gin.Context) {
		user := c.MustGet("user").(models.UserAuth)
		user.Password = ""
		c.JSON(200, user)
	})

	rg.PUT("/me", middleware.VerifyAuth(), func(c *gin.Context) {
		user := c.MustGet("user").(models.UserAuth)

		var body struct {
			Username string `json:"username" binding:"omitempty,min=3,max=50"`
			Email    string `json:"email" binding:"omitempty,email"`
			Password string `json:"password" binding:"omitempty,min=6"`
			FullName string `json:"fullName" binding:"omitempty"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"errors": err.Error()})
			return
		}

		if body.Username != "" && body.Username != user.Username {
			var existing models.UserAuth
			if err := config.DB.Where("username = ? AND id != ?", body.Username, user.ID).First(&existing).Error; err == nil {
				c.JSON(400, gin.H{"error": "Username already taken"})
				return
			}
			user.Username = body.Username
		}
		if body.Email != "" && body.Email != user.Email {
			var existing models.UserAuth
			if err := config.DB.Where("email = ? AND id != ?", body.Email, user.ID).First(&existing).Error; err == nil {
				c.JSON(400, gin.H{"error": "Email already taken"})
				return
			}
			user.Email = body.Email
		}
		if body.FullName != "" {
			user.FullName = body.FullName
		}
		if body.Password != "" {
			user.Password = body.Password
			user.HashPassword()
		}

		config.DB.Save(&user)
		user.Password = ""
		c.JSON(200, gin.H{"message": "Profile updated", "user": user})
	})

	rg.GET("/users", middleware.VerifyAdmin(), func(c *gin.Context) {
		var users []models.UserAuth
		config.DB.Order("created_at DESC").Find(&users)
		for i := range users {
			users[i].Password = ""
		}
		c.JSON(200, users)
	})

	rg.PUT("/users/:id/role", middleware.VerifyAdmin(), func(c *gin.Context) {
		var body struct {
			Role string `json:"role" binding:"required,oneof=admin manager viewer"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid role"})
			return
		}

		var user models.UserAuth
		if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		user.Role = body.Role
		config.DB.Save(&user)
		c.JSON(200, gin.H{"message": "Role updated", "user": gin.H{"id": user.ID, "role": user.Role}})
	})
}
