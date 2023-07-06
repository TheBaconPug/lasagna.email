package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func init() {
	LoadConfig()
	CreateMongoClient()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"Domains": Config.Domains, "Address": RandomString(10)})
	})

	router.GET("/inbox/:address", func(c *gin.Context) {
		c.HTML(http.StatusOK, "inbox.html", gin.H{"Address": c.Param("address")})
	})

	router.GET("/inbox/email/:id", func(c *gin.Context) {
		id := c.Param("id")

		email := GetEmailById(id)

		c.HTML(http.StatusOK, "email.html", gin.H{"Id": id, "Sender": email.Sender, "Subject": email.Subject, "Body": email.Body})
	})

	router.GET("/api/inbox/:address", func(c *gin.Context) {
		address := c.Param("address")

		emails := GetInbox(address)

		c.JSON(http.StatusOK, gin.H{
			"emails": emails,
		})
	})

	router.GET("/api/email/:id", func(c *gin.Context) {
		id := c.Param("id")

		email := GetEmailById(id)

		c.JSON(http.StatusOK, email)
	})

	router.POST("/api/callback", func(c *gin.Context) {
		c.Request.ParseForm()

		recipient := c.Request.PostForm.Get("recipient")
		sender := c.Request.PostForm.Get("sender")
		subject := c.Request.PostForm.Get("subject")

		var emailBody string

		bodyHtml := c.Request.PostForm.Get("body-html")
		if bodyHtml != "" {
			emailBody = bodyHtml
		} else {
			emailBody = c.Request.PostForm.Get("stripped-html")
		}

		if ValidateEmail(recipient) && ValidateEmail(sender) {
			guid := xid.New()

			CreateEmail(Email{Id: guid.String(), Recipient: recipient, Sender: sender, Subject: subject, Body: emailBody})

			log.Printf("received email | recipient: %s | sender: %s | id: %s\n", recipient, sender, guid.String())

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

	log.Println("Starting server on port " + Config.Port)

	router.Run(fmt.Sprintf(":%s", Config.Port))
}
