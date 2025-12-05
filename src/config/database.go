package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_URI")

	if os.Getenv("NODE_ENV") == "production" {
		dsn += " sslmode=require"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // mirrors `logging: false`
	})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// mirrors Sequelize's pool: { max: 10, min: 0, acquire: 30000, idle: 10000 }
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetConnMaxLifetime(0)

	DB = db
	fmt.Println("Database connected")
}
