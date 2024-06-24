package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

func InitDatabase() bool {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	dbName := os.Getenv("MONGODB_DBNAME")
	if dbName == "" {
		log.Println("The 'MONGODB_DBNAME' environmental variable is not set. Defaulting to 'WeebsKingdom'.")
		dbName = "WeebsKingdom"
	}
	MongoDB = client.Database(dbName)
	return true
}
