package env

import (
	"github.com/joho/godotenv"
)

// Load environment variables from file
func Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
