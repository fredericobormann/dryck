package models

import "github.com/jinzhu/gorm"

// User stores information about one user
type User struct {
	gorm.Model
	Name string
}

// getID returns the id of a user
func (u *User) getID() uint {
	return u.Model.ID
}
