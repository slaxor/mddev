package main

import (
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func fsWatch(filename string, wsCh chan string, done chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if filepath.Base(ev.Name) == filepath.Base(filename) {
					switch ev.Op {
					default:
						continue
					case fsnotify.Write:
						wsCh <- readMd(filename)
					}
				}

			case err := <-watcher.Errors:
				log.Fatalf("fserror: %v", err)
			}
		}
	}()

	err = watcher.Add(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}
	<-done
	watcher.Close()
}
