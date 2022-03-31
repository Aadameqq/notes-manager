package structs

type Note struct {
	Id      string `bson:"_id"`
	OwnerId string `bson:"ownerId"`
	Content string `bson:"content"`
}
