package initializers

import (
	"log"
	"os"
)

func CheckTLSFilesExistence() {
	// Check if the cert.pem and key.pem files exist
	var files = []string{"../tls/cert.pem", "../tls/key.pem"}

	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			log.Fatal(err)
		}
	}
}
