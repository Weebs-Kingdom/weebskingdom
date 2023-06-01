package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"weebskingdom/main/crypt"
	"weebskingdom/main/database"
	"weebskingdom/main/database/models"
)

//check if the header has the api key

func LoginToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get jwt token from header
		token, err := c.Cookie("auth")
		dataIgnore, exists := c.Get("ignoreAuth")
		var ignoreAuth bool
		if exists {
			ignoreAuth = dataIgnore.(bool)
		} else {
			ignoreAuth = false
		}

		if err != nil {
			fmt.Println("No token")
			if ignoreAuth {
				c.Next()
				return
			}
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		if token == "" {
			fmt.Println("No token")
			//Check if ignoreAuth is true and if it is, ignore the auth
			if ignoreAuth {
				c.Next()
				return
			}
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		jwt, err := crypt.ParseJwt(token)
		if err != nil {
			fmt.Println("Invalid token")
			if ignoreAuth {
				c.Next()
				return
			}
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		id, err := primitive.ObjectIDFromHex(jwt["userId"].(string))
		if err != nil {
			fmt.Println("Invalid token")
			if ignoreAuth {
				c.Next()
				return
			}
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		res := database.MongoDB.Collection("user").FindOne(c, bson.M{
			"_id": id,
		})

		if res.Err() != nil {
			fmt.Println("User not found")
			if ignoreAuth {
				c.Next()
				return
			}
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		//check time
		if jwt["exp"] != nil {
			if jwt["exp"].(float64) < jwt["iat"].(float64) {
				fmt.Println("Token expired")
				if ignoreAuth {
					c.Next()
					return
				}
				c.JSON(401, gin.H{
					"message": "Unauthorized",
				})
				c.Abort()
				return
			}
		}
		var dUser models.User
		err = res.Decode(&dUser)

		if err == nil {
			c.Set("user", dUser)
		}

		c.Set("userId", jwt["userId"])
		c.Set("loggedIn", true)
		c.Next()
	}
}
