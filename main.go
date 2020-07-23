package main

import (
	"fmt"
	"github.com/fredericobormann/dryck/db"
	"github.com/fredericobormann/dryck/handler"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	"os"
	"strconv"
	"time"
)

var dryckdb *db.DB

func main() {
	databasePassword := os.Getenv("POSTGRES_PASSWORD")
	dryckdb = db.New("postgres", "host=postgres user=postgres dbname=postgres password="+databasePassword+" sslmode=disable")

	defer dryckdb.Close()

	dryckhandler := handler.New(dryckdb)

	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsPrice": formatAsPrice,
		"formatAsTime":  formatAsTime,
	})
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", dryckhandler.HandleIndex)
	router.POST("/new-user", dryckhandler.HandleNewUser)
	router.GET("/user/:user_id", dryckhandler.HandleUserPage)
	router.POST("/purchase/:user_id", dryckhandler.HandlePurchase)
	router.POST("/delete-purchase/:user_id", dryckhandler.HandleDeletePurchase)

	router.POST("/new-payment/:user_id", dryckhandler.HandlePayment)
	router.POST("/delete-payment/:user_id", dryckhandler.HandleDeletePayment)

	router.Run()
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
