package models

import "github.com/jinzhu/gorm"

// Drink stores information about one type of drink
type Drink struct {
	gorm.Model
	Name  string
	Price int
}

// getID returns the ID of a drink
func (d *Drink) getID() uint {
	return d.Model.ID
}
