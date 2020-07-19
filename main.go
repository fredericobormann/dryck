package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

func main() {
	db, err := gorm.Open("postgres", "host=localhost user=postgres dbname=postgres password=mysecretpassword sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
}
