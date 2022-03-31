package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"web-app-go/structs"
	"web-app-go/utils"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Nickname string             `bson:"nickname" validate:"required,gte=5,lte=50"`
	Password string             `bson:"password" validate:"required,gte=6,lte=50"`
}

func LoginHandler(db *mongo.Database, jwtSecret string) func(serverContext *gin.Context) {
	return func(serverContext *gin.Context) {
		validate := validator.New()

		coll := db.Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var body structs.User

		var result structs.User
		serverContext.ShouldBindBodyWith(&body, binding.JSON)

		err := validate.Struct(body)

		if err != nil {
			serverContext.Status(http.StatusBadRequest)
			return
		}
		err = coll.FindOne(ctx, bson.D{{"nickname", body.Nickname}}).Decode(&result)

		if err == mongo.ErrNoDocuments {
			serverContext.Status(http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(body.Password))

		if err != nil {
			serverContext.Status(http.StatusUnauthorized)
			return
		}

		user := structs.User{result.Id, result.Nickname, ""}

		claims := &utils.Claims{
			Id:       user.Id.Hex(),
			Username: user.Nickname,
			RegisteredClaims: &jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 20)),
			},
		}

		tokenSigner := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

		token, err := tokenSigner.SignedString([]byte(jwtSecret))

		if err != nil {
			serverContext.Status(http.StatusInternalServerError)
			return
		}

		serverContext.JSON(http.StatusOK, gin.H{"token": token})
	}
}
