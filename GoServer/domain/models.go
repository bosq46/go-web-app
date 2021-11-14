package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Account  string
	Name     string
	Password string
}
