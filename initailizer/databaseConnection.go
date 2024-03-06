package initailizer

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ? 這邊特別新建全局變量就是讓DB可以在其他地方使用，而不會限制在ConnectToDatabase這個函數裡面
var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("fail to connect to database")
	}
}
