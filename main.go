package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Email struct {
	gorm.Model
	Id        string
	Recipient string
	Sender    string
	Subject   string
	Body      string
}

func main() {
	config := LoadConfig()

	db, err := gorm.Open(sqlite.Open("emails.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Email{})

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Static("/assets", "./assets")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"Domains": config.Domains, "Address": RandomString(10)})
	})

	router.GET("/inbox/:email", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inbox.html", gin.H{"Address": c.Param("email")})
	})

	router.GET("/inbox/email/:id", func(c *gin.Context) {
		var emails []Email

		id := c.Param("id")

		db.Where("id = ?", id).Find(&emails)

		email := emails[0]

		c.HTML(http.StatusOK, "email.html", gin.H{"Id": id, "Sender": email.Sender, "Subject": email.Subject, "Body": email.Body})
	})

	router.GET("/api/inbox/:email", func(c *gin.Context) {
		var emails []Email

		email := c.Param("email")

		db.Where("recipient = ?", email).Find(&emails)

		c.JSON(http.StatusOK, gin.H{
			"emails": emails,
		})
	})

	router.GET("/api/email/:id", func(c *gin.Context) {
		var emails []Email

		id := c.Param("id")

		db.Where("id = ?", id).Find(&emails)

		c.JSON(http.StatusOK, gin.H{
			"emails": emails,
		})
	})

	router.POST("/api/callback", func(c *gin.Context) {
		var emailBody string

		c.Request.ParseForm()

		fmt.Println(c.Request.PostForm)

		recipient := c.Request.PostForm["recipient"][0]
		sender := c.Request.PostForm["sender"][0]
		subject := c.Request.PostForm["subject"][0]

		if body, exists := c.Request.PostForm["body-html"]; exists {
			emailBody = body[0]
		} else {
			emailBody = c.Request.PostForm["stripped-html"][0]
		}

		if ValidateEmail(recipient) && ValidateEmail(sender) {
			guid := xid.New()

			db.Create(&Email{Id: guid.String(), Recipient: recipient, Sender: sender, Subject: subject, Body: emailBody})

			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})

			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid email address",
		})
	})

	router.Run(fmt.Sprintf(":%s", config.Port))
}
