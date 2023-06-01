package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"weebskingdom/main/api/middleware"
	"weebskingdom/main/crypt"
	"weebskingdom/main/database"
	"weebskingdom/main/database/models"
)

func InitApi(r *gin.Engine) {
	apiAuth := r.Group("/api/dev")

	userAuth := apiAuth.Group("/api/user")
	userAuth.Use(middleware.LoginToken())

	initApis(apiAuth)
	initUserApi(userAuth)
}

func initUserApi(r *gin.RouterGroup) {
	r.GET("/auth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Logged in",
		})
	})
}

func initApis(r *gin.RouterGroup) {
	type Register struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	r.POST("/register", func(c *gin.Context) {
		var register Register
		err := c.ShouldBindBodyWith(&register, binding.JSON)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid JSON",
			})
			return
		}

		fmt.Println(register.Email)
		email := database.MongoDB.Collection("user").FindOne(c, bson.M{
			"email": strings.ToLower(register.Email),
		})

		if email.Err() == nil {
			c.JSON(400, gin.H{
				"message": "E-Mail already taken",
			})
			return
		} else {
			if register.Email == "" {
				c.JSON(400, gin.H{
					"message": "Email cannot be empty",
				})
				return
			}

			if register.Password == "" {
				c.JSON(400, gin.H{
					"message": "Password cannot be empty",
				})
				return
			}

			if len(register.Password) < 8 {
				c.JSON(400, gin.H{
					"message": "Password must be at least 8 characters",
				})
				return
			}

			if len(register.Password) > 64 {
				c.JSON(400, gin.H{
					"message": "Password must be at most 64 characters",
				})
				return
			}

			//check if password contains at least one uppercase letter and one number and one special character and one lowercase letter
			if !strings.ContainsAny(register.Password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
				c.JSON(400, gin.H{
					"message": "Password must contain at least one uppercase letter",
				})
				return
			}

			if !strings.ContainsAny(register.Password, "abcdefghijklmnopqrstuvwxyz") {
				c.JSON(400, gin.H{
					"message": "Password must contain at least one lowercase letter",
				})
				return
			}

			if !strings.ContainsAny(register.Password, "0123456789") {
				c.JSON(400, gin.H{
					"message": "Password must contain at least one number",
				})
				return
			}

			if !strings.ContainsAny(register.Password, "!@#$%^&*()-_+={}[]:;\"'<>,.?/|") {
				c.JSON(400, gin.H{
					"message": "Password must contain at least one special character",
				})
				return
			}

			password, err := crypt.HashPassword(register.Password)
			if err != nil {
				log.Println("Error while hashing password")
				log.Println(err)
				c.JSON(500, gin.H{
					"message": "Internal Server Error",
				})
				return
			}
			_, err2 := database.MongoDB.Collection("user").InsertOne(c, models.User{
				ID:                primitive.NewObjectID(),
				Email:             strings.ToLower(register.Email),
				Password:          password,
				IsVerifiedDiscord: false,
				IsAdmin:           false,
				IsDeveloper:       false,
			})

			if err2 != nil {
				log.Println("Error while creating user")
				log.Println(err2)
				c.JSON(500, gin.H{
					"message": "Internal Server Error",
				})
				return
			} else {
				c.JSON(200, gin.H{
					"message": "User created successfully",
					"status":  200,
				})
			}
		}
	})

	r.POST("/login", func(c *gin.Context) {
		var login Register
		err := c.ShouldBindBodyWith(&login, binding.JSON)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid JSON",
			})
			return
		}
		var dUser models.User
		user := database.MongoDB.Collection("user").FindOne(c, bson.M{
			"email": strings.ToLower(login.Email),
		})

		if user.Err() != nil {
			c.JSON(400, gin.H{
				"message": "Invalid username or password",
			})
		} else {
			err := user.Decode(&dUser)
			if err != nil {
				c.JSON(500, gin.H{
					"message": "Internal Server Error",
				})
				return
			}

			if crypt.CheckPasswordHash(login.Password, dUser.Password) {
				token, err := crypt.GenerateLoginToken(dUser.ID)
				if err != nil {
					c.JSON(500, gin.H{
						"message": "Internal Server Error",
					})
					return
				}

				c.JSON(200, gin.H{
					"message": "Logged in",
					"token":   token,
					"status":  200,
				})
			} else {
				c.JSON(400, gin.H{
					"message": "Invalid username or password",
				})
			}
		}
	})
}
