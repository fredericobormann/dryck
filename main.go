package main

import (
	"dryck/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
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
	router.SetFuncMap(template.FuncMap{
		"formatAsPrice": formatAsPrice,
		"formatAsTime":  formatAsTime,
	})
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", handleIndex)
	router.POST("/new-user", handleNewUser)
	router.GET("/user/:user_id", handleUserPage)
	router.POST("/purchase/:user_id", handlePurchase)

	router.Run()
}

func getAllUsers() []models.User {
	var allUsers []models.User
	db.Find(&allUsers)
	return allUsers
}

func getUsername(userId uint) string {
	var user models.User
	db.Where("id = ?", userId).First(&user)
	return user.Name
}

func createNewUser(username string) {
	newUser := models.User{Name: username}
	db.Create(&newUser)
}

func getAllDrinks() []models.Drink {
	var allDrinks []models.Drink
	db.Find(&allDrinks)
	return allDrinks
}

func getPurchasesOfUser(userId uint) []models.Purchase {
	var purchases []models.Purchase
	db.Preload("Product").Where("customer_id = ?", userId).Order("purchase_time desc").Find(&purchases)
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

func handleNewUser(c *gin.Context) {
	newUserName := c.PostForm("new-user-name")
	createNewUser(newUserName)

	c.Redirect(http.StatusMovedPermanently, "/")
}

func handleUserPage(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	username := getUsername(uint(userId))
	purchases := getPurchasesOfUser(uint(userId))
	drinks := getAllDrinks()

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"title":     username,
			"username":  username,
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

func formatAsPrice(cents int) string {
	if cents%100 >= 10 {
		return strconv.FormatInt(int64(cents/100), 10) + "," + strconv.FormatInt(int64(cents%100), 10) + "€"
	} else {
		return strconv.FormatInt(int64(cents/100), 10) + ",0" + strconv.FormatInt(int64(cents%100), 10) + "€"
	}
}

func formatAsTime(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%02d.%02d.%d", day, month, year)
}
