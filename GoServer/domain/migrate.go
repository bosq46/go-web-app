package domain

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbName = "data.sqlite3"

func Migrate() {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	// d, _ := db.DB()
	// defer d.Close()
	db.AutoMigrate(&User{})
}
