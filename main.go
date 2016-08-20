package main

import (
	"flag"
	"log"
)

func main() {
	file := flag.String("f", "", "-f <filename.md>")
	addr := flag.String("a", "127.0.0.1:4050", "-a <ip:port> # default: 127.0.0.1:4050")
	flag.Parse()
	log.Printf("watching %v", *file)
	log.Printf("now open  http://%v", *addr)
	done := make(chan bool)
	wsCh := make(chan string)
	go fsWatch(*file, wsCh, done)
	startHttp(*addr, *file, wsCh)
}
