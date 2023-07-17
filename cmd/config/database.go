package config

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var logger = logrus.New()

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:admin@tcp(127.0.0.1:3306)/go_gRPC_server_products?parseTime=true"), &gorm.Config{})
	if err != nil {
		logger.Println("Failed to Connect Database")
	}

	logger.Println("Database Connected")

	return db
}
