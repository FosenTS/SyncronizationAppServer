package watcher

import (
	"GolangSync/server"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func StartWatch(config configWatcher) {
	// setup watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	// use goroutine to start the watcher
	go func() {
		for {
			select {
			// provide the list of events to monitor
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					updateSync(event.Name, config.mask)
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					updateSync(event.Name, config.mask)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					updateSync(event.Name, config.mask)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					updateSync(event.Name, config.mask)
				}
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					updateSync(event.Name, config.mask)
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	// provide the directory to monitor
	err = watcher.Add(config.pathDirectory)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func updateSync(pathFile string, maskArg string) {
	fileExtension := filepath.Ext(pathFile)
	if fileExtension == maskArg {
		fmt.Println("File is edit: " + pathFile)
		server.SendAllClientMessage(server.Clients, pathFile)
	}
}

type configWatcher struct {
	pathDirectory string
	mask          string
}

func NewConfigWatcher() *configWatcher {
	fmt.Println("Uses directory: " + os.Args[1])
	fmt.Println("File mask use: " + os.Args[2])
	return &configWatcher{
		pathDirectory: os.Args[1],
		mask:          os.Args[2],
	}
}
