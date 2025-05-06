package main

import (
	"rest-api/app/config"
	"rest-api/app/db"
	"rest-api/app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080") //localhost:8080
}
