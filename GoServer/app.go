package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"forest.work/m/domain"

	"github.com/gorilla/sessions"
	_ "gorm.io/driver/sqlite"
)

const sesKey = "go-server-app-session-key"
const sesLoginKey = "go-server-app-session-key-login"
const templateDir = "templates/user/"
const port = ":8080"

var cs *sessions.CookieStore = sessions.NewCookieStore([]byte(sesKey))

func noTemplate() *template.Template {
	src := "<html><body><h1>NO TEMPLATE.</h1></body></html>"
	tmp, _ := template.New("index").Parse(src)
	return tmp
}

func setLogin(r *http.Request, w http.ResponseWriter, login bool, name string) {
	ses, _ := cs.Get(r, sesLoginKey)
	ses.Values["login"] = login
	ses.Values["name"] = name
	ses.Save(r, w)
}

func getLogin(r *http.Request) (string, bool) {
	ses, _ := cs.Get(r, sesLoginKey)
	login, exist := ses.Values["login"]
	name := ses.Values["name"]
	if !exist {
		ses.Values["login"] = false
		ses.Values["name"] = ""
	}
	return name.(string), login.(bool)
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
		// TODO: 毒抜き

		// save DB
		result, err := domain.RegisterUser(name, password)
		fmt.Printf("result %s\n", strconv.FormatBool(result))
		fmt.Printf("err %s\n", err)
		// save DB
		if result && err == nil {
			setLogin(r, w, true, name)
			http.Redirect(w, r, "home", http.StatusFound)
		} else {
			msg = "登録に失敗:" + err.Error()
		}
	}

	item := struct {
		Title           string
		Message         string
		PostURL         string
		ResponseMessage string
	}{
		Title:           "新規ユーザ作成",
		Message:         "ユーザ名とパスワードを入力してください",
		PostURL:         "register",
		ResponseMessage: msg,
	}
	err = template.Execute(w, item)
	if err != nil {
		fmt.Println(err)
	}
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

	name, isLogin := getLogin(r)
	msg := ""
	// ログイン済みの場合
	if isLogin {
		msg = "login as " + name
		http.Redirect(w, r, "home", http.StatusFound)
	}

	// ログイン情報の受付
	if r.Method == "POST" {
		name := r.PostFormValue("name")
		password := r.PostFormValue("password")
		// TODO: 毒抜き
		//reflect.ValueOf(user).IsNil()
		if domain.LoginUser(name, password) {
			fmt.Println("login success.")
			setLogin(r, w, true, name)
			http.Redirect(w, r, "home", http.StatusFound)
		} else {
			fmt.Println("login failed.")
			msg = "パスワードが異なります．"
		}
	}

	// ログインフォームの表示
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
	setLogin(r, w, false, "")
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

	_, isLogin := getLogin(r)
	// ログイン済みの場合
	if !isLogin {
		http.Redirect(w, r, "login", http.StatusFound)
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
	http.HandleFunc("/register", Register)
	http.HandleFunc("/home", UserList)
	http.HandleFunc("/", UserList)

	fmt.Println("Running on http://127.0.0.1" + port + "/ (Press CTRL+C to quit)")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
