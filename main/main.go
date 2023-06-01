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

	crypt.InitJwt()
	database.InitDatabase()

	api.InitApi(r)
	core.LoadTemplates(r)
	core.LoadServerAssets(r)

	//set address
	address := ""
	if os.Getenv("ADDRESS") != "" {
		address = os.Getenv("ADDRESS")
	}
	r.Run(address)
}
