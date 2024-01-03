package handlers

import (
	"net/http"
	"html/template"
)

func RenderHTMLTemplate(w http.ResponseWriter, templateName string, data ...interface{}) {
	tmpl, err := template.ParseFiles("internal/templates/" + templateName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if data is provided
	if len(data) > 0 {
		err = tmpl.Execute(w, data[0])
	} else {
		err = tmpl.Execute(w, nil)
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
