package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/russross/blackfriday"
)

func readMd(filename string) string {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	html := blackfriday.MarkdownCommon(buf)
	return fmt.Sprintf("%s", html)
	return fmt.Sprintf("%s", html)
}
