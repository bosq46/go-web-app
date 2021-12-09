package domain

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(name string, pass string) bool {
	user, err := FindUser(name)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(pass))
	fmt.Println(err)
	return err == nil
}

func RegisterUser(name, password string) (bool, error) {
	user, err := FindUser(name)
	fmt.Println("Registered = name:" + user.Name + " pass: " + string(user.Password))

	// TODO: Validation
	// nil はすでに登録済みを表すのでエラーを返す
	if err == nil {
		return false, errors.New("the name is already registered")
	}
	CreateUser(name, password)
	return true, nil
}

func generatePassword(pass string) []byte {
	// 2^12 = 4096ビットでハッシュ作成
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		fmt.Println(err)
	}
	return hash
}

func DeleteUser(id int) (bool, error) {
	return DeleteUserRecord(id)
	// error handling
}
