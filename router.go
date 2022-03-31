package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"web-app-go/handlers"
	"web-app-go/middlewares"
	"web-app-go/utils"
)

func SetRouters(router *gin.RouterGroup, db *mongo.Database, cfg utils.Config) {
	requireAuthGroup := router.Group("/")
	requireAuthGroup.Use(middlewares.AuthorizeHandler(cfg.JWT_SECRET))
	{
		requireAuthGroup.POST("/notes", handlers.CreateNoteHandler(db))
		requireAuthGroup.GET("/notes", handlers.ReadAllNotesHandler(db))
		requireAuthGroup.GET("/notes/:id", handlers.ReadOneNoteHandler(db))
		requireAuthGroup.DELETE("/notes/:id", handlers.DeleteNoteHandler(db))
	}
	router.POST("/login", handlers.LoginHandler(db, cfg.JWT_SECRET))

	router.POST("/register", handlers.RegisterHandler(db))

}
