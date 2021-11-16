package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
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

func FindUser(name string, password string) User {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	var user User
	db.Where("name = ? AND password = ?", name, generatePassword(password)).Find(&user)
	return user
}

func generatePassword(pass string) []byte {
	// 2^12 = 4096ビットでハッシュ作成
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		fmt.Println(err)
	}
	return hash
}
