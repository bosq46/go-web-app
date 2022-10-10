package sqlite3

import (
	"errors"
	"fmt"
	"strconv"

	"go-web-app/adapter/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbName = "data.sqlite3"

type User struct {
	gorm.Model
	Name           string `gorm:"unique;not null"`
	HashedPassword []byte `gorm:"not null"`
}

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
	hashedPass := utils.GeneratePassword(password)
	db.Create(&User{Name: name, HashedPassword: hashedPass})
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

func FindUserOnUnscoped(name string) (User, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	var user User
	err = db.Unscoped().Where("name = ?", name).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, err
	}
	return user, nil
}

func FindUserById(id int) (User, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	var user User
	err = db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, err
	}
	return user, nil
}

func DeleteUserRecord(id int) (bool, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return false, errors.New("DB can not open")
	}
	user, err := FindUserById(id)
	if err != nil {
		fmt.Println(err)
		return false, errors.New("User ID not found: " + strconv.Itoa(id))
	}
	db.Where("id = ?", user.ID).Delete(&User{})
	return true, nil
}

func UpdateUserRecord(id int, name string, encryptedPassword []byte) (bool, error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return false, errors.New("DB can not open")
	}
	user, err := FindUserById(id)
	if err != nil {
		fmt.Println(err)
		return false, errors.New("User ID not found: " + strconv.Itoa(id))
	}
	// db.Model(&user).Omit("id").Updates(map[string]interface{}{"name": name, "hashedPassword": encryptedPassword})
	db.Model(&user).Updates(User{
		Name:           name,
		HashedPassword: encryptedPassword,
	})
	return true, nil
}
