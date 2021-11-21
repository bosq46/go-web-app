package domain

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbName = "data.sqlite3"

func Migrate() {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	// d, _ := db.DB()
	// defer d.Close()
	db.AutoMigrate(&User{})
}

func CreateUser(name string, password string) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	db.Create(&User{Name: name, Password: generatePassword(password)})
}

func ListUser() []User {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	var users = []User{}
	db.Find(&users)
	return users
}

func FindUser(name string) (User, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	var user User
	err = db.Where("name = ?", name).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, err
	}
	return user, nil
}
