package middleware

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"weebskingdom/main/crypt"
	"weebskingdom/main/database"
)

//check if the header has the api key

func LoginToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get jwt token from header
		token, err := c.Cookie("auth")
		if err != nil {
			val, exists := c.Get("ignoreAuth")
			if exists {
				if val.(bool) {
					c.Next()
					return
				}
			}
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		if token == "" {
			//Check if ignoreAuth is true and if it is, ignore the auth
			val, exists := c.Get("ignoreAuth")
			if exists {
				if val.(bool) {
					c.Next()
					return
				}
			}
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		jwt, err := crypt.ParseJwt(token)
		if err != nil {
			val, exists := c.Get("ignoreAuth")
			if exists {
				if val.(bool) {
					c.Next()
					return
				}
			}
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		res := database.MongoDB.Collection("users").FindOne(c, bson.M{
			"_id": jwt["userId"],
		})

		if res == nil {
			val, exists := c.Get("ignoreAuth")
			if exists {
				if val.(bool) {
					c.Next()
					return
				}
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
				val, exists := c.Get("ignoreAuth")
				if exists {
					if val.(bool) {
						c.Next()
						return
					}
				}
				c.JSON(401, gin.H{
					"message": "Unauthorized",
				})
				c.Abort()
				return
			}
		}

		c.Set("userId", jwt["userId"])
		c.Set("user", res)
		c.Set("loggedIn", true)
		c.Next()
	}
}
