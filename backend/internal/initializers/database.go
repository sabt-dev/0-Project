package initializers

import (
	"log"
	"os"
	"time"

	"github.com/sabt-dev/0-Project/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
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
	// Open a connection to the SQLite database using modernc.org/sqlite
	DB, err = gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        config.AppConfig.DBHost,
	}, &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
}
