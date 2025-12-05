package routes

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(rg *gin.RouterGroup) {
	cld, _ := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)

	rg.POST("", func(c *gin.Context) {
		file, err := c.FormFile("profilePhoto")
		if err != nil {
			c.JSON(400, gin.H{"error": "No file uploaded"})
			return
		}
		opened, _ := file.Open()
		defer opened.Close()

		result, err := cld.Upload.Upload(context.Background(), opened, uploader.UploadParams{
			Folder: "mahbere-maryam-profiles",
		})
		if err != nil {
			c.JSON(400, gin.H{"error": "Upload failed"})
			return
		}
		c.JSON(200, gin.H{"fileUrl": result.SecureURL, "filename": result.PublicID})
	})
}
