package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(pass string) []byte {
	// 2^12 = 4096ビットでハッシュ作成
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		fmt.Println(err)
	}
	return hash
}
