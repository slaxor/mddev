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
	// htmlFlags := 0
	// htmlFlags |= blackfriday.HTML_TOC
	// htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	// htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	// htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	// renderer := blackfriday.HtmlRenderer(htmlFlags, filename, "/mdapp.css")
	// extensions := 0
	// extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	// extensions |= blackfriday.EXTENSION_TABLES
	// extensions |= blackfriday.EXTENSION_FENCED_CODE
	// extensions |= blackfriday.EXTENSION_AUTOLINK
	// extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	// extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	// html := blackfriday.Markdown(buf, renderer, extensions)
	// fmt.Printf("%+v", buf)
	// input := os.Stdin.Read(&buf)
	html := blackfriday.MarkdownCommon(buf)
	return fmt.Sprintf("%s", html)
	return fmt.Sprintf("%s", html)
}
