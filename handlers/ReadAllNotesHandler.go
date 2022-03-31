package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
	"web-app-go/structs"
)

func ReadAllNotesHandler(db *mongo.Database) func(ctx *gin.Context) {
	return func(serverContext *gin.Context) {
		coll := db.Collection("notes")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stringOwnerId := fmt.Sprintf("%v", serverContext.MustGet("id"))

		var results []structs.Note

		cursor, err := coll.Find(ctx, bson.D{{"ownerId", stringOwnerId}})

		if err = cursor.All(ctx, &results); err != nil {
			serverContext.Status(http.StatusNotFound)
		}

		serverContext.JSON(200, results)
	}
}
