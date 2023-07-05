package main

import (
	"GolangSync/server"
	"GolangSync/watcher"
)

func main() {
	configWatcher := watcher.NewConfigWatcher()
	go watcher.StartWatch(*configWatcher)
	server.StartWebsocket()
}
