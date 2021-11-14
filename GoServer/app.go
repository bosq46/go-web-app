package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"forest.work/m/domain"

	"github.com/gorilla/sessions"
	_ "gorm.io/driver/sqlite"
)

var sesKey = "go-server-app-session-key"
var cs *sessions.CookieStore = sessions.NewCookieStore([]byte(sesKey))

func showParams(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Printf("r.Form: %s\n", r.Form)
	fmt.Fprintf(w, "r.Form: %s\n", r.Form)

	fmt.Printf("r.URL.Path: %s\n", r.URL.Path)
	fmt.Fprintf(w, "r.URL.Path: %s\n", r.URL.Path)

	fmt.Printf("r.URL.Scheme: %s\n", r.URL.Scheme)
	fmt.Fprintf(w, "r.URL.Scheme: %s\n", r.URL.Scheme)

	fmt.Printf("r.Form[\"url_long\"]: %s\n", r.Form["url_long"])
	fmt.Fprintf(w, "r.Form[\"url_long\"]: %s\n", r.Form["url_long"])

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Fprintf(w, "key:%s\n", k)
		fmt.Println("val:", v)
		fmt.Fprintf(w, "val:%s\n", strings.Join(v, ""))
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	tmp, template_err := template.ParseFiles("../templates/sample.html")
	if template_err != nil {
		log.Fatal(template_err)
		panic(template_err)
	}

	item := struct {
		Title   string
		Message string
		Items   []string
	}{
		Title:   "sample page",
		Message: "this is sample. <br>",
		Items:   []string{"one", "two", "tree"},
	}

	err := tmp.Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("../templates/post.html")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	msg := "Please Input."
	if r.Method == "POST" {
		name := r.PostFormValue("name")
		password := r.PostFormValue("password")
		msg = "name:" + name + " pass:" + password
	}
	item := struct {
		Title   string
		Message string
	}{
		Title:   "Send values",
		Message: msg,
	}

	err = template.Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

func notemp() *template.Template {
	src := "<html><body><h1>NO TEMPLATE.</h1></body></html>"
	tmp, _ := template.New("index").Parse(src)
	return tmp
}

func login(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles(
		"../templates/login.html",
		"../templates/header.html",
		"../templates/footer.html",
	)
	if err != nil {
		log.Fatal(err)
		// panic(err)
		template = notemp()
	}
	ses, _ := cs.Get(r, "sample-session")

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
	}{
		Title:   "Session",
		Message: msg,
		Account: name,
	}
	err = template.Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	domain.Migrate()
	http.HandleFunc("/params", showParams)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/post", post)
	http.HandleFunc("/login", login)

	port := ":8080"
	fmt.Println("Running on http://127.0.0.1" + port + "/ (Press CTRL+C to quit)")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
