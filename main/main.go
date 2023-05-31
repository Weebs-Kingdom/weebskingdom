package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"weebskingdom/main/api"
	"weebskingdom/main/core"
	"weebskingdom/main/crypt"
	"weebskingdom/main/database"
)

func main() {
	r := gin.Default()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	crypt.InitJwt()
	database.InitDatabase()

	api.InitApi(r)
	core.LoadTemplates(r)
	core.LoadServerAssets(r)

	r.Run()
}
