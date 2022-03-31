package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"web-app-go/structs"
)

func RegisterHandler(db *mongo.Database) func(*gin.Context) {
	return func(serverContext *gin.Context) {
		coll := db.Collection("users")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var body structs.User

		serverContext.ShouldBindBodyWith(&body, binding.JSON)
		validate := validator.New()
		err := validate.Struct(body)

		if err != nil {
			serverContext.Status(http.StatusBadRequest)
			return
		}

		var result structs.User

		err = coll.FindOne(ctx, bson.D{{"nickname", body.Nickname}}).Decode(&result)

		if err == nil {
			serverContext.Status(http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

		if err != nil {
			serverContext.Status(http.StatusInternalServerError)
			return
		}

		_, err = coll.InsertOne(ctx, bson.D{{"nickname", body.Nickname}, {"password", string(hashedPassword)}})

		if err != nil {
			serverContext.Status(http.StatusInternalServerError)
			return
		}

		serverContext.Status(http.StatusCreated)

	}
}
