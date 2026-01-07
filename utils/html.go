package utils

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func RenderHtml(w http.ResponseWriter, name string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	tmpl := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/"+name,
		),
	)
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
