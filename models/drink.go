package models

import "github.com/jinzhu/gorm"

type Drink struct {
	gorm.Model
	Name  string
	Price int
}

func (d *Drink) getId() uint {
	return d.Model.ID
}
