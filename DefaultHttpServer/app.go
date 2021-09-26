package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

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

func main() {
	http.HandleFunc("/", showParams)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
