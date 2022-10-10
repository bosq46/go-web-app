package controller

import (
	"errors"
	"fmt"
	"go-web-app/adapter/sqlite3"
	"go-web-app/adapter/utils"

	"golang.org/x/crypto/bcrypt"
)

func PostUsers(name, password string) (bool, error) {
	_, err := sqlite3.FindUserOnUnscoped(name)

	// TODO: Validation
	// nil はすでに登録済みを表すのでエラーを返す
	if err == nil {
		return false, errors.New("the name is already registered")
	}
	sqlite3.CreateUser(name, password)
	return true, nil
}

func LoginUser(name string, pass string) bool {
	user, err := sqlite3.FindUser(name)
	if err != nil {
		return false
	}
	fmt.Println(user.HashedPassword)
	fmt.Println(pass)
	fmt.Printf("%T", user.HashedPassword)
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(pass))
	fmt.Println(err)
	return err == nil
}

func UpdateUser(id int, name string, rawPassword string) (bool, error) {
	pw := utils.GeneratePassword(rawPassword)
	sqlite3.UpdateUserRecord(id, name, pw)
	return true, nil
}
