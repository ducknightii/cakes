package main

import "net/http"

func main() {
	uri := "baidu.com"
	http.Get(uri)
}
