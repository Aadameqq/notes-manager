package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
	"web-app-go/structs"
)

type readOneNoteHandlerRequest struct {
	Id string `json:"id" validate:"required"`
}

func ReadOneNoteHandler(db *mongo.Database) func(ctx *gin.Context) {
	return func(serverContext *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		idFromQuery := serverContext.Param("id")

		params := readOneNoteHandlerRequest{idFromQuery}

		validate := validator.New()

		err := validate.Struct(params)

		if err != nil {
			serverContext.Status(http.StatusBadRequest)
			return
		}

		var result structs.Note

		id, err := primitive.ObjectIDFromHex(params.Id)

		if err != nil {
			serverContext.Status(http.StatusInternalServerError)
			return
		}

		ownerId := serverContext.MustGet("id")

		err = db.Collection("notes").FindOne(ctx, &bson.D{{"_id", id}, {"ownerId", fmt.Sprintf("%v", ownerId)}}).Decode(&result)

		if err != nil {
			serverContext.Status(http.StatusNotFound)
			return
		}

		serverContext.JSON(http.StatusOK, result)
	}
}
