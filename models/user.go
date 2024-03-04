package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID   `bson:"_id"`
	Name      string               `bson:"name"`
	Email     string               `bson:"email"`
	IsPremium bool                 `bson:"is_premium,omitempty"`
	Todos     []primitive.ObjectID `bson:"todos,omitempty"`
	Password  string               `bson:"password"`
	CreatedAt primitive.DateTime   `bson:"created_at,omitempty"`
}

type RegisterResponse struct {
	UserID primitive.ObjectID
	Token  string
}
