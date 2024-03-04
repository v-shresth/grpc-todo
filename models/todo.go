package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      primitive.ObjectID `bson:"user_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description,omitempty"`
	Status      bool               `bson:"status,omitempty"`
	Priority    string             `bson:"priority,omitempty"`
	CreateTime  primitive.DateTime `bson:"create_time,omitempty"`
	UpdateTime  primitive.DateTime `bson:"update_time,omitempty"`
	DeadLine    primitive.DateTime `bson:"deadline,omitempty"`
}

type ListTodoFilter struct {
	Limit int32
	Page  int32
}

type ListTodoRes struct {
	Todos []Todo
	Count int64
}
