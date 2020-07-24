package main

import (
	"github.com/fredericobormann/dryck/db"
	"github.com/fredericobormann/dryck/format"
	"github.com/fredericobormann/dryck/handler"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"html/template"
	"os"
)

func main() {
	databasePassword := os.Getenv("POSTGRES_PASSWORD")
	dryckdb := db.New("postgres", "host=postgres user=postgres dbname=postgres password="+databasePassword+" sslmode=disable")

	defer dryckdb.Close()

	dryckhandler := handler.New(dryckdb)

	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsPrice": format.AsPrice,
		"formatAsTime":  format.AsTime,
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
