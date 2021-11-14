package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forest.work/m/domain"

	"github.com/gorilla/sessions"
	_ "gorm.io/driver/sqlite"
)

const sesKey = "go-server-app-session-key"
const sesLoginKey = "go-server-app-session-key-login"
const templateDir = "templates/"
const port = ":8080"

var cs *sessions.CookieStore = sessions.NewCookieStore([]byte(sesKey))

func NoTemplate() *template.Template {
	src := "<html><body><h1>NO TEMPLATE.</h1></body></html>"
	tmp, _ := template.New("index").Parse(src)
	return tmp
}

func Login(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		templateDir+"login.html",
		templateDir+"header.html",
		templateDir+"footer.html",
	)
	if err != nil {
		fmt.Println("Can no find template file.")
		template = NoTemplate()
	}
	ses, _ := cs.Get(r, sesLoginKey)

	if r.Method == "POST" {
		ses.Values["login"] = nil
		ses.Values["name"] = nil
		nm := r.PostFormValue("name")
		pw := r.PostFormValue("password")
		if nm == pw {
			ses.Values["login"] = true
			ses.Values["name"] = nm
		}
		ses.Save(r, w)
	}
	isLogin, _ := ses.Values["login"].(bool)
	name, _ := ses.Values["name"].(string)
	msg := "no login"
	if isLogin {
		msg = "login as " + name
	}
	item := struct {
		Title   string
		Message string
		Account string
		PostURL string
	}{
		Title:   "Session",
		Message: msg,
		Account: name,
		PostURL: "test/post",
	}
	err = template.Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		templateDir+"new.html",
		templateDir+"header.html",
		templateDir+"footer.html",
	)
	if err != nil {
		fmt.Println("Can no find template file.")
		template = NoTemplate()
	}

	ses, _ := cs.Get(r, sesLoginKey)
	if r.Method == "POST" {
		ses.Values["login"] = nil
		ses.Values["name"] = nil
		nm := r.PostFormValue("name")
		pw := r.PostFormValue("password")

		// save DB
		domain.CreateUser(nm, pw)
		// save DB
		if nm == pw {
			ses.Values["login"] = true
			ses.Values["name"] = nm
		}
		ses.Save(r, w)
		http.Redirect(w, r, "users", http.StatusFound)
	}

	item := struct {
		Title   string
		Message string
		PostURL string
	}{
		Title:   "新規ユーザ作成",
		Message: "ユーザ名とパスワードを入力してください",
		PostURL: "new",
	}
	err = template.Execute(w, item)
	if err != nil {
		fmt.Println(err)
	}
}

func UserList(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		templateDir+"users.html",
		templateDir+"header.html",
		templateDir+"footer.html",
	)
	if err != nil {
		fmt.Println("Can no find template file.")
		template = NoTemplate()
	}

	ses, _ := cs.Get(r, sesLoginKey)
	if r.Method == "POST" {
		ses.Values["login"] = nil
		ses.Values["name"] = nil
		nm := r.PostFormValue("name")
		pw := r.PostFormValue("password")

		// save DB
		domain.CreateUser(nm, pw)
		// save DB
		if nm == pw {
			ses.Values["login"] = true
			ses.Values["name"] = nm
		}
		ses.Save(r, w)
		http.Redirect(w, r, "users", http.StatusFound)
	}

	var userNames []string
	userList := domain.ListUser()
	for _, user := range userList {
		userNames = append(userNames, user.Name)
	}
	item := struct {
		Title     string
		Message   string
		UserNames []string
	}{
		Title:     "My Home",
		Message:   "ユーザの編集が可能です",
		UserNames: userNames,
	}
	err = template.Execute(w, item)
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	domain.Migrate()
	// Test Page
	http.HandleFunc("/test/params", ShowParams)
	http.HandleFunc("/test/hello", Hello)
	http.HandleFunc("/test/post", Post)
	// User Login Page
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/new", NewUser)
	http.HandleFunc("/users", UserList)

	fmt.Println("Running on http://127.0.0.1" + port + "/ (Press CTRL+C to quit)")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
