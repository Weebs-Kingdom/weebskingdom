package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"time"
	"weebskingdom/api/middleware"
	"weebskingdom/crypt"
	"weebskingdom/database"
	"weebskingdom/database/models"
)

func InitApi(r *gin.Engine) {
	apiAuth := r.Group("/api/dev")

	userAuth := r.Group("/api/user")
	userAuth.Use(middleware.LoginToken())

	adminAuth := r.Group("/api/admin")
	adminAuth.Use(middleware.LoginToken())
	adminAuth.Use(middleware.VerifyAdmin())

	initApis(apiAuth)
	initUserApi(userAuth)
	initAdminApi(adminAuth)
	initSickUrls(r)
}

func initSickUrls(r *gin.Engine) {
	discordUrl := "https://discord.gg/gK8F7kuybv"
	r.GET("/discord", func(c *gin.Context) {
		c.Redirect(302, discordUrl)
	})

	r.GET("/disc", func(c *gin.Context) {
		c.Redirect(302, discordUrl)
	})
}

func initUserApi(r *gin.RouterGroup) {
	r.GET("/auth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Logged in",
		})
	})
}

func initAdminApi(r *gin.RouterGroup) {
	r.GET("/auth", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Logged in",
		})
	})

	r.DELETE("/contact", func(c *gin.Context) {
		type Contact struct {
			ID string `json:"contactID"`
		}

		var contact Contact
		err := c.ShouldBindBodyWith(&contact, binding.JSON)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid JSON",
			})
			return
		}

		dataId, err := primitive.ObjectIDFromHex(contact.ID)

		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid ID",
			})
			return
		}

		_, err = database.MongoDB.Collection("contact").DeleteOne(c, bson.M{
			"_id": dataId,
		})
		if err != nil {
			c.JSON(403, gin.H{
				"message": "This contact doesn't exist",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Deleted contact",
		})
	})
}

func initApis(r *gin.RouterGroup) {
	type Register struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type ContactForm struct {
		Email   string `json:"email"`
		Subject string `json:"subject"`
		Message string `json:"message"`
		Topic   string `json:"topic"`
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

	r.POST("/contact", func(c *gin.Context) {
		var contact ContactForm
		err := c.ShouldBindBodyWith(&contact, binding.JSON)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid JSON",
			})
			return
		}

		//get user
		c.Set("ignoreAuth", true)
		middleware.LoginToken()(c)
		c.Set("ignoreAuth", false)

		isLoggedIn := c.GetBool("loggedIn")
		var userId = primitive.NilObjectID
		if isLoggedIn {
			dUser, ok := c.Get("user")
			if ok {
				user := dUser.(models.User)
				userId = user.ID
			}
		}

		_, err = database.MongoDB.Collection("contact").InsertOne(c, models.Contact{
			ID:         primitive.NewObjectID(),
			User:       userId,
			Message:    contact.Message,
			Email:      contact.Email,
			Subject:    contact.Subject,
			Topic:      contact.Topic,
			DateIssued: primitive.NewDateTimeFromTime(time.Now()),
		})
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Internal Server Error",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Contact sent",
			"status":  200,
		})

	})
}
