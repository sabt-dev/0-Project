package initializers

import (
	"github.com/sabt-dev/0-Project/internal/models"
)

func SyncDatabase() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		return
	}
}
