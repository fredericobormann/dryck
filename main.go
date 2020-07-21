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
	httpPassword, httpBasicAuthActive := os.LookupEnv("HTTP_PASSWORD")

	defer dryckdb.Close()

	dryckhandler := handler.New(dryckdb)

	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsPrice": format.AsPrice,
		"formatAsTime":  format.AsTime,
	})
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	var authorized *gin.RouterGroup
	if httpBasicAuthActive {
		authorized = router.Group("/", gin.BasicAuth(gin.Accounts{
			"dryck": httpPassword,
		}))
	} else {
		authorized = router.Group("/")
	}

	authorized.GET("/", dryckhandler.HandleIndex)
	authorized.POST("/new-user", dryckhandler.HandleNewUser)
	authorized.GET("/user/:user_id", dryckhandler.HandleUserPage)
	authorized.POST("/purchase/:user_id", dryckhandler.HandlePurchase)
	authorized.POST("/delete-purchase/:user_id", dryckhandler.HandleDeletePurchase)

	authorized.POST("/new-payment/:user_id", dryckhandler.HandlePayment)
	authorized.POST("/delete-payment/:user_id", dryckhandler.HandleDeletePayment)

	router.Run()
}
