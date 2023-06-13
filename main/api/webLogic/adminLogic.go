package webLogic

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"weebskingdom/database"
	"weebskingdom/database/models"
)

type AdminContact struct {
	Contacts []ContactStuff `json:"contacts" bson:"contacts"`
}

type ContactStuff struct {
	ContactID  string      `json:"contactID" bson:"_id"`
	Email      string      `json:"email" bson:"email"`
	Subject    string      `json:"subject" bson:"subject"`
	Message    string      `json:"message" bson:"message"`
	Topic      string      `json:"topic" bson:"topic"`
	User       models.User `json:"user" bson:"user"`
	DateIssued string      `json:"dateIssues" bson:"dateIssues"`
	FoundUser  bool        `json:"foundUser" bson:"foundUser"`
}

func adminContact(c *gin.Context) any {
	var contacts []models.Contact
	var contactStuff []ContactStuff
	cursor, err := database.MongoDB.Collection("contact").Find(c, bson.M{})
	if err != nil {
		return AdminContact{
			Contacts: []ContactStuff{},
		}
	}
	err = cursor.All(c, &contacts)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		return AdminContact{
			Contacts: []ContactStuff{},
		}
	}
	for i := range contacts {
		var user models.User

		foundUser := false
		if contacts[i].User != primitive.NilObjectID {
			err = database.MongoDB.Collection("user").FindOne(c, bson.M{
				"_id": contacts[i].User,
			}).Decode(&user)
			foundUser = true
		}

		contactStuff = append(contactStuff, ContactStuff{
			ContactID:  contacts[i].ID.Hex(),
			Email:      contacts[i].Email,
			Subject:    contacts[i].Subject,
			Message:    contacts[i].Message,
			Topic:      contacts[i].Topic,
			User:       user,
			DateIssued: contacts[i].DateIssued.Time().Format("01.02.2006 15:04:05"),
			FoundUser:  foundUser,
		})
	}
	return AdminContact{
		Contacts: contactStuff,
	}
}
