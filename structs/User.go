package structs

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Nickname string             `bson:"nickname"`
	Password string             `bson:"password"`
}
