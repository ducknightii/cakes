package main

import (
	"fmt"
	"io/ioutil"

	"github.com/h2non/filetype"
)

func main() {
	buf, err := ioutil.ReadFile("test.png")
	fmt.Println(len(buf), err)

	kind, _ := filetype.Match(buf)
	fmt.Println(filetype.Image(buf))
	if kind == filetype.Unknown {
		fmt.Println("Unknown file type")
		return
	}

	fmt.Printf("File type: %s. MIME: %+v\n", kind.Extension, kind.MIME) // File type: png. MIME: {Type:image Subtype:png Value:image/png}

}
