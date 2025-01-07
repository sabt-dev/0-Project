package initializers

import (
	"log"

	"github.com/sabt-dev/0-Project/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := config.AppConfig.DBUser + ":" + config.AppConfig.DBPassword + "@tcp(" + config.AppConfig.DBHost + ":" + config.AppConfig.DBPort + ")/" + config.AppConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}
