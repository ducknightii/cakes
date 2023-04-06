package main

import (
	"fmt"
	"net/url"
	"strings"
)

func main() {
	u := "postgres://user:pass@host.com:5432/path?k=v&c=d#f"

	urlParse, _ := url.Parse(u)

	cuts := strings.Split(urlParse.Path, "/")

	fmt.Printf("urlParse.Path:%#v -> %s\n", urlParse.Path, cuts[len(cuts)-1])
	fmt.Printf("urlParse.RawPath:%#v\n", urlParse.RawPath)
	fmt.Printf("%#v\n", urlParse)

}
