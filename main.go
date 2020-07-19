package main

import (
	"dryck/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"strconv"
	"time"
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

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handleIndex)
	router.GET("/user/:user_id", handleUserPage)
	router.POST("/purchase/:user_id", handlePurchase)

	router.Run()
}

func getAllUsers() []models.User {
	var allUsers []models.User
	db.Find(&allUsers)
	return allUsers
}

func getAllDrinks() []models.Drink {
	var allDrinks []models.Drink
	db.Find(&allDrinks)
	return allDrinks
}

func getPurchasesOfUser(userId uint) []models.Purchase {
	var purchases []models.Purchase
	db.Preload("Product").Where("customer_id = ?", userId).Find(&purchases)
	return purchases
}

func purchaseDrink(userId uint, drinkId uint) {
	purchase := models.Purchase{CustomerID: userId, ProductID: drinkId, PurchaseTime: time.Now()}
	db.Create(&purchase)
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
			"users": getAllUsers(),
		},
	)
}

func handleUserPage(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	purchases := getPurchasesOfUser(uint(userId))
	drinks := getAllDrinks()

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"title":     "Users",
			"userId":    userId,
			"drinks":    drinks,
			"purchases": purchases,
		},
	)
}

func handlePurchase(c *gin.Context) {
	drinkId, _ := strconv.ParseUint(c.PostForm("drink"), 10, 64)
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	purchaseDrink(uint(userId), uint(drinkId))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}
