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
				log.Printf("%v == %v", filepath.Base(ev.Name), filename)
				if filepath.Base(ev.Name) == filepath.Base(filename) {
					switch ev.Op {
					case fsnotify.Create:
						log.Printf("fsevent: Created %v", ev.Name)
					case fsnotify.Write:
						wsCh <- readMd(filename)
						log.Printf("fsevent: Written %v", ev.Name)
					case fsnotify.Remove:
						log.Printf("fsevent: Deleted %v", ev.Name)
					case fsnotify.Rename:
						log.Printf("fsevent: Renamed %v", ev.Name)
					case fsnotify.Chmod:
						log.Printf("fsevent: Changed permissions for %v", ev.Name)

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

	/* ... do stuff ... */
	watcher.Close()
}
