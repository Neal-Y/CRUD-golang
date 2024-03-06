package main

import (
	"crud_go_gin/controllers"
	"crud_go_gin/initailizer"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	controllers.SetupPostsRouter(r)
	controllers.SetupAuthRouter(r)

	r.Run()
}

func init() {
	initailizer.LoadLocalVariables()
	initailizer.ConnectToDatabase()
}
