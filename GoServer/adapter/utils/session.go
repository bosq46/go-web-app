package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const sesKey = "go-server-app-session-key"
const sesLoginKey = "go-server-app-session-key-login"

var cs *sessions.CookieStore = sessions.NewCookieStore([]byte(sesKey))

func SetLogin(r *http.Request, w http.ResponseWriter, login bool, name string) {
	// session 初期設定
	cs.Options.HttpOnly = true
	// cs.Options.Secure = true
	ses, _ := cs.Get(r, sesLoginKey)
	ses.Values["login"] = login
	ses.Values["name"] = name
	ses.Save(r, w)
}

func GetLogin(r *http.Request, w http.ResponseWriter) (string, bool) {
	ses, _ := cs.Get(r, sesLoginKey)
	login, exist := ses.Values["login"]
	name := ses.Values["name"]
	if !exist {
		ses.Values["login"] = false
		ses.Values["name"] = ""
		ses.Save(r, w)
		return "", false
	}
	return name.(string), login.(bool)
}
