package main

import "GolangSync/watcher"

func main() {
	configWatcher := watcher.NewConfigWatcher()
	watcher.StartWatch(*configWatcher)
}
