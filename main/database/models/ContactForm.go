package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Contact struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	User       primitive.ObjectID `json:"user" bson:"user"`
	Email      string             `json:"email" bson:"email"`
	Subject    string             `json:"subject" bson:"subject"`
	Message    string             `json:"message" bson:"message"`
	Topic      string             `json:"topic" bson:"topic"`
	DateIssued primitive.DateTime `json:"dateIssues" bson:"dateIssues"`
}
