package controller

import (
	"errors"
	"fmt"
	"go-web-app/adapter/sqlite3"
	"go-web-app/adapter/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

const templateDir = "templates/user/"

func noTemplate() *template.Template {
	src := "<html><body><h1>NO TEMPLATE.</h1></body></html>"
	tmp, _ := template.New("index").Parse(src)
	return tmp
}

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

func Register(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		templateDir+"register.html",
		templateDir+"header.html",
		templateDir+"footer.html",
	)
	if err != nil {
		fmt.Println("Can no find template file.")
		template = noTemplate()
	}

	msg := ""
	if r.Method == "POST" {
		name := r.PostFormValue("name")
		password := r.PostFormValue("password")

		result, err := PostUsers(name, password)
		if result && err == nil {
			utils.SetLogin(r, w, true, name)
			http.Redirect(w, r, "home", http.StatusFound)
		} else {
			msg = "登録に失敗:" + err.Error()
		}
	}

	item := struct {
		Title           string
		PostURL         string
		ResponseMessage string
	}{
		Title:           "新規ユーザ作成ページ",
		PostURL:         "register",
		ResponseMessage: msg,
	}
	err = template.Execute(w, item)
	if err != nil {
		fmt.Println(err)
	}
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

func Login(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		templateDir+"login.html",
		templateDir+"header.html",
		templateDir+"footer.html",
	)
	if err != nil {
		fmt.Println("Can no find template file.")
		template = noTemplate()
	}

	name, isLogin := utils.GetLogin(r, w)
	msg := ""
	// ログイン済みの場合
	if isLogin {
		msg = "ログイン済み: " + name
		http.Redirect(w, r, "home", http.StatusFound)
	}

	// ログイン情報の受付
	if r.Method == "POST" {
		name := r.PostFormValue("name")
		password := r.PostFormValue("password")

		PostUsers(name, password)
		if LoginUser(name, password) {
			fmt.Println("login success.")
			utils.SetLogin(r, w, true, name)
			http.Redirect(w, r, "home", http.StatusFound)
		} else {
			fmt.Println("login failed.")
			msg = "パスワードが異なります"
		}
	}

	// ログインフォームの表示
	item := struct {
		Title           string
		PostURL         string
		ResponseMessage string
	}{
		Title:           "ログインフォーム",
		PostURL:         "login",
		ResponseMessage: msg,
	}
	err = template.Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	utils.SetLogin(r, w, false, "")
	http.Redirect(w, r, "login", http.StatusFound)
}

func UserList(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		templateDir+"users.html",
		templateDir+"header.html",
		templateDir+"footer.html",
	)
	if err != nil {
		fmt.Println("Can no find template file.")
		template = noTemplate()
	}

	_, isLogin := utils.GetLogin(r, w)
	// ログイン済みの場合
	if !isLogin {
		http.Redirect(w, r, "login", http.StatusFound)
	}

	UserInfo := make(map[uint]string)
	for _, user := range sqlite3.ListUser() {
		// fmt.Printf("user info -> %d : %s\n", user.ID, user.Name)
		UserInfo[user.ID] = user.Name
	}
	item := struct {
		Title    string
		UserInfo map[uint]string
	}{
		Title:    "ホーム画面",
		UserInfo: UserInfo,
	}
	err = template.Execute(w, item)
	if err != nil {
		fmt.Println(err)
	}
}

func UserEdit(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		templateDir+"edit.html",
		templateDir+"header.html",
		templateDir+"footer.html",
	)
	if err != nil {
		fmt.Println("Can no find template file.")
		template = noTemplate()
	}

	_, isLogin := utils.GetLogin(r, w)
	// ログインしていない場合
	if !isLogin {
		http.Redirect(w, r, "login", http.StatusFound)
	}

	// 編集の受付
	if r.Method == "POST" {
		method := r.PostFormValue("method")
		if method == "delete" {
			userId := r.PostFormValue("id")
			fmt.Println("Delete ->" + userId)
			userIdInt, _ := strconv.Atoi(userId)
			sqlite3.DeleteUserRecord(userIdInt)
		} else if method == "put" {
			targetStrID := r.PostFormValue("id")
			targetID, targetIDErr := strconv.Atoi(targetStrID)
			if targetIDErr != nil {
				fmt.Println("Invalid user ID")
			}
			targetName := r.PostFormValue("name")
			targetPassword := r.PostFormValue("password")
			fmt.Printf("Update -> %d / %s / %s\n", targetID, targetName, targetPassword)

			UpdateUser(targetID, targetName, targetPassword)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}

	userID := r.FormValue("user")
	item := struct {
		Title  string
		UserID string
	}{
		Title:  "ユーザ編集",
		UserID: userID,
	}
	err = template.Execute(w, item)
	if err != nil {
		fmt.Println(err)
	}
}
