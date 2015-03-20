package web

import (
	// "appengine"
	// "appengine/user"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/api", newHandler)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Rawr")
}
