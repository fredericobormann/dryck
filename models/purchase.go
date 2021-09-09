package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Purchase stores information about one purchase
type Purchase struct {
	gorm.Model
	Customer     User
	CustomerID   uint
	Product      Drink
	ProductID    uint
	PurchaseTime time.Time
	Price        int
}
