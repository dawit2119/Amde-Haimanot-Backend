package main

import (
	"AmdeHaimanot-Backend-Go/src/config"
	"AmdeHaimanot-Backend-Go/src/middleware"
	"AmdeHaimanot-Backend-Go/src/models"
	"AmdeHaimanot-Backend-Go/src/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not loaded:", err)
	} else {
		log.Println("✅ .env loaded successfully")
	}
	config.ConnectDatabase()

	if os.Getenv("NODE_ENV") != "production" {
		if err := models.Migrate(config.DB); err != nil {
			log.Fatal("migration failed:", err)
		}
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/uploads", "./uploads")

	router.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Backend is running"}) })
	router.GET("/api/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "OK"}) })

	api := router.Group("/api")
	{
		routes.RegisterAuthRoutes(api.Group("/auth"))
		routes.RegisterUserRoutes(api.Group("/users"))
		routes.RegisterImportRoutes(api.Group("/users/import")) // register before /users to avoid :id capturing "import"
		routes.RegisterCategoryRoutes(api.Group("/categories"))
		routes.RegisterUploadRoutes(api.Group("/upload"))
		// routes.RegisterAnalyticsRoutes(api.Group("/analytics"))
		// routes.RegisterSectionRoutes(api.Group("/sections"))
		// routes.RegisterCourseRoutes(api.Group("/courses"))
		routes.RegisterCourseRegistrationRoutes(api.Group("/course-registrations"))
		routes.RegisterAttendanceRoutes(api.Group("/attendance"))
		routes.RegisterQuestionnaireRoutes(api.Group("/questionnaires"))
	}
	router.Use(middleware.ErrorHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("🚀 Server running on port " + port)
	router.Run(":" + port)
}
