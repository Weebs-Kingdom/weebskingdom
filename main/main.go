package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"weebskingdom/api"
	"weebskingdom/core"
	"weebskingdom/crypt"
	"weebskingdom/database"
	"weebskingdom/env"
)

//go:embed web
var Files embed.FS

func main() {
	_, err := Files.ReadDir("web")
	if err != nil {
		log.Println("Failed to read public files - this is likely a problem during compilation. Exiting...")
		return
	}

	env.Files = Files

	r := gin.Default()
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	//Banner
	log.Println("\n" + env.BANNER + "\nWeebs Kingdom website" + "\nVersion: " + env.VERSION)

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
