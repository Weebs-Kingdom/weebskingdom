package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"weebskingdom/api"
	"weebskingdom/core"
	"weebskingdom/crypt"
	"weebskingdom/database"
)

func main() {
	r := gin.Default()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	crypt.KeySetup()
	database.InitDatabase()

	api.InitApi(r)
	core.LoadTemplates(r)
	core.LoadServerAssets(r)

	//set address
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	r.Run(":" + port)
}
