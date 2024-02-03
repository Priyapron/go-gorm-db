// models/Register.go
package models

import "gorm.io/gorm"

type Register struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Hash     string
}
