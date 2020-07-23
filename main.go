package main

import (
	"fmt"
	"github.com/fredericobormann/dryck/db"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

var dryckdb *db.DB

func main() {
	databasePassword := os.Getenv("POSTGRES_PASSWORD")
	dryckdb = db.New("postgres", "host=postgres user=postgres dbname=postgres password="+databasePassword+" sslmode=disable")

	defer dryckdb.Close()

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
	router.POST("/delete-purchase/:user_id", handleDeletePurchase)

	router.POST("/new-payment/:user_id", handlePayment)
	router.POST("/delete-payment/:user_id", handleDeletePayment)

	router.Run()
}

// Handles requests to the index page
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
			"users": dryckdb.GetAllUsers(),
		},
	)
}

// Handles creation of a new user
func handleNewUser(c *gin.Context) {
	newUserName := c.PostForm("new-user-name")
	dryckdb.CreateNewUser(newUserName)

	c.Redirect(http.StatusMovedPermanently, "/")
}

// Handles user page with purchase history and the option to buy new drinks
func handleUserPage(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	username := dryckdb.GetUsername(uint(userId))
	totalDebt := dryckdb.GetTotalDebtOfUser(uint(userId))
	purchases := dryckdb.GetPurchasesOfUser(uint(userId))
	payments := dryckdb.GetAllPaymentsOfUser(uint(userId))
	drinks := dryckdb.GetAllDrinks()

	c.HTML(
		http.StatusOK,
		"user.html",
		gin.H{
			"title":     username,
			"username":  username,
			"userId":    userId,
			"totalDebt": totalDebt,
			"drinks":    drinks,
			"purchases": purchases,
			"payments":  payments,
		},
	)
}

// Handles a new purchase
func handlePurchase(c *gin.Context) {
	drinkId, _ := strconv.ParseUint(c.PostForm("drink"), 10, 64)
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	dryckdb.PurchaseDrink(uint(userId), uint(drinkId))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// Handles the deletion of a purchase
func handleDeletePurchase(c *gin.Context) {
	purchaseId, _ := strconv.ParseUint(c.PostForm("delete_purchase"), 10, 64)
	dryckdb.DeletePurchase(uint(purchaseId))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// Handles a new payment
func handlePayment(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	paymentAmount, _ := strconv.ParseInt(c.PostForm("payment_amount"), 10, 64)

	dryckdb.AddPayment(uint(userId), int(paymentAmount))
	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

func handleDeletePayment(c *gin.Context) {
	paymentId, _ := strconv.ParseUint(c.PostForm("delete_payment"), 10, 64)
	dryckdb.DeletePayment(uint(paymentId))

	c.Redirect(http.StatusMovedPermanently, "/user/"+c.Param("user_id"))
}

// Formats a given cent amount to Eurostring
func formatAsPrice(cents int) string {
	result := ""
	posCents := cents
	if cents < 0 {
		posCents = -cents
	}
	if posCents%100 >= 10 {
		result = strconv.FormatInt(int64(posCents/100), 10) + "," + strconv.FormatInt(int64(posCents%100), 10) + "€"
	} else {
		result = strconv.FormatInt(int64(posCents/100), 10) + ",0" + strconv.FormatInt(int64(posCents%100), 10) + "€"
	}
	if cents < 0 {
		return "-" + result
	}
	return result
}

// Formats a timestamp so it's human readable
func formatAsTime(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%02d.%02d.%d", day, month, year)
}
