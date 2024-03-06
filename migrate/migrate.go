package main

import (
	"crud_go_gin/initailizer"
	"crud_go_gin/models"
)

func init() {
	initailizer.LoadLocalVariables()
	initailizer.ConnectToDatabase()
}

func main() {
	initailizer.DB.AutoMigrate(&models.Post{})
	initailizer.DB.AutoMigrate(&models.User{})
}
