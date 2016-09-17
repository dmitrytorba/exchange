package main

import (
	"net/http"
	"text/template"
)

func executeTemplate(w http.ResponseWriter, name string, status int, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	return tpls[name].ExecuteTemplate(w, "base", data)
}

var tpls = map[string]*template.Template{
	"signup": newTemplate("templates/base.html", "templates/signup.html"),
	"home":     newTemplate("templates/base.html", "templates/home.html"),
	"error":    newTemplate("templates/base.html", "templates/error.html"),
	"settings": newTemplate("templates/base.html", "templates/settings.html"),
}

func newTemplate(files ...string) *template.Template {
	return template.Must(template.New("*").ParseFiles(files...))
}
