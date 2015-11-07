package hello
import (
	"net/http"
)

func init() {
	http.HandleFunc("/rest", root)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "static/index.html")
	w.WriteHeader(http.StatusFound)
}

