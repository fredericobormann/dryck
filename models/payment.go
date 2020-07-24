package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// Payment stores information about one payment
type Payment struct {
	gorm.Model
	User        User
	UserID      uint
	PaymentTime time.Time
	Amount      int
}
