package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Alert struct {
	ID         primitive.ObjectID `bson:"_id"`
	KafkaTopic QueueTopic         `bson:"type"`
	Type       string             `bson:"sub_type"`
	UserID     primitive.ObjectID `bson:"user_id"`
	TodoID     primitive.ObjectID `bson:"todo_id"`
}
