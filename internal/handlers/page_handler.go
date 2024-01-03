package handlers

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    RenderHTMLTemplate(w, "index.html")
}
