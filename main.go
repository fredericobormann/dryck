package main

import (
	"dryck/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost user=postgres dbname=postgres password=mysecretpassword sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Drink{})
	db.AutoMigrate(&models.Purchase{})

	var users []models.User
	db.Find(&users)
	fmt.Println(users)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handleIndex)

	router.Run()
}

func getAllUsers(db *gorm.DB) []models.User {
	var allUsers []models.User
	db.Find(&allUsers)
	return allUsers
}

func handleIndex(c *gin.Context) {
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"title": "Home Page",
			"users": getAllUsers(db),
		},
	)
}
