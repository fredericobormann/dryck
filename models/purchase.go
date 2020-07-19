package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Purchase struct {
	gorm.Model
	Customer     User
	UserID       uint
	Product      Drink
	ProductID    uint
	PurchaseTime time.Time
}
