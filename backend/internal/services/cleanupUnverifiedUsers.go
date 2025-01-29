package services

import (
	"log"
	"time"

	"github.com/sabt-dev/0-Project/internal/initializers"
	"github.com/sabt-dev/0-Project/internal/models"
)

func CleanupUnverifiedUsers() {
    for {
        time.Sleep(time.Minute) // Run the cleanup job every minute

        // Calculate the cutoff time
        cutoffTime := time.Now().Add(-10 * time.Minute)

        // Delete unverified users who were created more than 10 minutes ago
        result := initializers.DB.Where("verified = ? AND created_at < ?", false, cutoffTime).Delete(&models.User{})
        if result.Error != nil || result.RowsAffected != 0 {
            return
        }
		if result.RowsAffected > 0 {
            log.Printf("Deleted %d unverified users", result.RowsAffected)
            return
		}
    }
}