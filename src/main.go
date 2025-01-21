package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	app := gin.Default()
	app.POST("/install", runInstallation)

	if address := os.Getenv("ADDRESS"); address != "" {
		port := os.Getenv("PORT")
		app.Run(address + ":" + port)
	}
	app.Run()
}
