package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logger = logrus.New()

func ConnectDB() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		logger.Printf("failed load env file %s", err.Error())
	}
	
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		logger.Printf("failed connect to database %s", err.Error())
	}

	logger.Println("database connected")

	return db
}
