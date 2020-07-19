package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name string
}

func (u *User) getId() uint {
	return u.Model.ID
}
