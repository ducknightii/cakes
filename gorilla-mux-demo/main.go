package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthz)
	http.Handle("/", r)

}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
