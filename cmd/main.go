package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/moikk-app/database"
)

func main() {
	//start sv
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	gin.SetMode(gin.ReleaseMode)
	database.InitDB()
	StartServer()
}
