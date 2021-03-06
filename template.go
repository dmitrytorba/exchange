package main

import (
	"log"
	"net/http"
	"text/template"
)

func createTemplateHandler(name string, isPublic bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//user.go
		user, err := getUserFromCookie(r)
		if user == nil {
			if isPublic {
				user = &User{}
			} else {
				http.Redirect(w, r, "/login", 302)
				return
			}
		}
		err = executeTemplate(w, name, 200, map[string]interface{}{
			"User": user,
		}) 
		if err != nil {
			executeTemplate(w, "error", 500, map[string]interface{}{
				"Error": err.Error(),
			})
			log.Println(err)
		}
	})
}


func executeTemplate(w http.ResponseWriter, name string, status int, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	return tpls[name].ExecuteTemplate(w, "base", data)
}

var tpls = map[string]*template.Template{
	"home":     newTemplate("templates/base.html", "templates/home.html"),
	"error":    newTemplate("templates/base.html", "templates/error.html"),
	"settings": newTemplate("templates/base.html", "templates/settings.html"),
	"signup":   newTemplate("templates/base.html", "templates/signup.html"),
	"login":    newTemplate("templates/base.html", "templates/login.html"),
}

func newTemplate(files ...string) *template.Template {
	return template.Must(template.New("*").ParseFiles(files...))
}
