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
        
        tx := initializers.DB.Begin()
        if tx.Error != nil {
            log.Printf("Failed to start transaction: %v", tx.Error)
            return
        }

        // Delete unverified users who were created more than 10 minutes ago
        result := tx.Where("verified = ? AND created_at < ?", false, cutoffTime).Delete(&models.User{})
        if result.Error != nil || result.RowsAffected != 0 {
            tx.Rollback()
            if tx.Error != nil {
                log.Printf("Failed to rollback transaction: %v", tx.Error)
                return
            }
        }
		if result.RowsAffected > 0 {
            if err := tx.Commit().Error; err != nil {
                tx.Rollback()
                if tx.Error != nil {
                    log.Printf("Failed to rollback transaction: %v", tx.Error)
                    return
                }
                log.Printf("Failed to commit transaction: %v", err)
                return
            }
			log.Printf("Deleted %d unverified users", result.RowsAffected)
		}
    }
}