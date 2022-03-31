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

type deleteNoteHandlerParams struct {
	Id string `json:"id" validate:"required"`
}

func DeleteNoteHandler(db *mongo.Database) func(ctx *gin.Context) {
	return func(serverContext *gin.Context) {
		notesColl := db.Collection("notes")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stringId := serverContext.Param("id")

		params := deleteNoteHandlerParams{Id: stringId}

		validate := validator.New()
		err := validate.Struct(params)

		if err != nil {
			serverContext.Status(http.StatusBadRequest)
			return
		}

		ownerId := serverContext.MustGet("id")

		note := structs.Note{Id: params.Id, OwnerId: fmt.Sprintf("%v", ownerId)}

		var result structs.Note

		objId, err := primitive.ObjectIDFromHex(note.Id)

		if err != nil {
			serverContext.Status(http.StatusInternalServerError)
			return
		}

		err = notesColl.FindOneAndDelete(ctx, &bson.D{{"_id", objId}, {"ownerId", note.OwnerId}}).Decode(&result)

		if err != nil {
			serverContext.Status(http.StatusNotFound)
			return
		}

		serverContext.Status(http.StatusOK)
	}
}
