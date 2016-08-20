package main

import "flag"

func main() {
	file := flag.String("f", "", "-f <filename.md>")
	addr := flag.String("a", "127.0.0.1:4050", "-a <ip:port> # default: 127.0.0.1:4050")
	flag.Parse()
	done := make(chan bool)
	wsCh := make(chan string)
	go fsWatch(*file, wsCh, done)
	startHttp(*addr, *file, wsCh)
}
