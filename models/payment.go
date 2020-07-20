package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	User        User
	UserID      uint
	PaymentTime time.Time
	Amount      int
}
