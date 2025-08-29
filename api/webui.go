package api

import "net/http"

func webUiHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/web.html")
}

func eulerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/cytoscape-euler.js")
}
