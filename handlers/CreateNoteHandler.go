package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
	"web-app-go/structs"
)

type createNoteHandlerBody struct {
	Content string `json:"content" validate:"required,gte=10,lte=1000"`
}

func CreateNoteHandler(db *mongo.Database) func(ctx *gin.Context) {
	return func(serverContext *gin.Context) {
		usersColl := db.Collection("users")
		notesColl := db.Collection("notes")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var body createNoteHandlerBody

		serverContext.ShouldBindBodyWith(&body, binding.JSON)
		validate := validator.New()
		err := validate.Struct(body)

		if err != nil {
			serverContext.Status(http.StatusBadRequest)
			return
		}

		var findUserResult structs.User

		stringOwnerId := serverContext.MustGet("id")

		var ownerId primitive.ObjectID

		ownerId, err = primitive.ObjectIDFromHex(fmt.Sprintf("%v", stringOwnerId))

		if err != nil {
			serverContext.Status(http.StatusInternalServerError)
			return
		}

		err = usersColl.FindOne(ctx, bson.D{{"_id", ownerId}}).Decode(&findUserResult)

		if err != nil {
			serverContext.Status(http.StatusBadRequest)
			return
		}

		note := structs.Note{OwnerId: findUserResult.Id.Hex(), Content: body.Content}

		_, err = notesColl.InsertOne(ctx, bson.D{{"ownerId", note.OwnerId}, {"content", note.Content}})

		if err != nil {
			serverContext.Status(http.StatusInternalServerError)
			return
		}

		serverContext.Status(http.StatusCreated)

	}
}
