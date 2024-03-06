package initailizer

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadLocalVariables() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
