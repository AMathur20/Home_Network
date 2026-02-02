package topology

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func WatchTopology(path string, callback func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Printf("Topology file modified: %s, reloading...", event.Name)
					callback()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error watching topology file:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Printf("Warning: Could not watch topology file %s: %v", path, err)
	} else {
		log.Printf("Watching topology file: %s", path)
	}
}
