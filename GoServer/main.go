package main

import (
	"fmt"
	"log"
	"net/http"

	"go-web-app/adapter/controller"
	"go-web-app/adapter/sqlite3"

	_ "gorm.io/driver/sqlite"
)

const port = ":8080"

func main() {
	sqlite3.Migrate()
	// Test Page
	http.HandleFunc("/test/params", controller.ShowParams)
	http.HandleFunc("/test/hello", controller.Hello)
	http.HandleFunc("/test/post", controller.Post)
	// User Login Page
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/logout", controller.Logout)
	http.HandleFunc("/register", controller.Register)
	http.HandleFunc("/edit", controller.UserEdit)
	http.HandleFunc("/", controller.UserList)
	fmt.Println("Running on http://127.0.0.1" + port + "/ (Press CTRL+C to quit)")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
