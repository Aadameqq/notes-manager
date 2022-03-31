package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"os"
	"time"
	"web-app-go/utils"
)

func ConnectDatabase() (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))

	if err != nil {
		panic(err)
	}

	var disconnect = func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}

	db := client.Database("notesappgo")
	fmt.Println("Connected to db!")
	return db, disconnect
}

func loadConfig() (cfg utils.Config) {
	f, err := os.Open("config.yaml")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	return
}

func main() {
	router := gin.Default()
	db, disconnect := ConnectDatabase()
	defer disconnect()

	routerGroup := router.Group("/api/v1/")

	cfg := loadConfig()
	SetRouters(routerGroup, db, cfg)
	router.Run("127.0.0.1:8080")

}
