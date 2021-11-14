package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func ShowParams(w http.ResponseWriter, r *http.Request) {
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

func Hello(w http.ResponseWriter, r *http.Request) {
	tmp, template_err := template.ParseFiles("templates/sample.html")
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

func Post(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("templates/post.html")
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
		PostURL string
	}{
		Title:   "Send values",
		Message: msg,
		PostURL: "post",
	}

	err = template.Execute(w, item)
	if err != nil {
		log.Fatal(err)
	}
}
