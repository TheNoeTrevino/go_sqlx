package main

import (
	"gotutorial/controllers"
	"gotutorial/db_client"

	"github.com/gin-gonic/gin"
)

func main() {
	db_client.InitializeDBConnection()

	router := gin.Default()

	router.POST("/", controllers.CreatePost)
	router.GET("/", controllers.GetPosts)
	router.GET("/:id", controllers.GetPost)

	err := router.Run(":5000")
	if err != nil {
		panic(err.Error())
	}
}
