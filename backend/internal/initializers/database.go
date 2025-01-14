package initializers

import (
	"log"
	"os"
	"time"

	"github.com/sabt-dev/0-Project/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
    // Configure the logger to be silent
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),    // io writer
        logger.Config{
            SlowThreshold:             time.Second,   // Slow SQL threshold
            LogLevel:                  logger.Warn,	  // Log level
            IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
            Colorful:                  true,          // Enable color
        },
    )
	dsn := config.AppConfig.DBUser + ":" + config.AppConfig.DBPassword + "@tcp(" + config.AppConfig.DBHost + ":" + config.AppConfig.DBPort + ")/" + config.AppConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
}
