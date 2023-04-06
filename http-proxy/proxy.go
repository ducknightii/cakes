package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	http.HandleFunc("/hello", proxyHello)
	http.ListenAndServe(":8091", nil)
}

func proxyHello(w http.ResponseWriter, r *http.Request) {
	u, _ := url.Parse("http://127.0.0.1:8090/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("http: proxy error!!!: %v\n", err)
		w.WriteHeader(http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}
