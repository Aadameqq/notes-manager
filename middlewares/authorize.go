package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"web-app-go/utils"
)

func AuthorizeHandler(jwtSecret string) func(*gin.Context) {
	return func(context *gin.Context) {

		headerValue := context.GetHeader("Authorization")
		if headerValue == "" {
			context.AbortWithStatus(http.StatusUnauthorized)
		}

		token, err := jwt.ParseWithClaims(headerValue, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Invalid signature")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || token.Valid != true {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*utils.Claims)

		context.Set("nickname", claims.Username)
		context.Set("id", claims.Id)
		context.Next()
	}
}
