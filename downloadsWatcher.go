// This script allows to track the downloads folder for changes.
package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// watch for folder changes
	if err := watcher.Add("/home/mastodilu/Downloads"); err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			fmt.Printf("%v %v\n", event.Op.String(), event.Name)

		// watch for errors
		case err := <-watcher.Errors:
			log.Println(err)
		}
	}

	/*
		scaricare dal web qualcosa genera questi eventi:
			CREATE /home/mastodilu/Downloads/fiore_azzurro_1024x768.jpg
			WRITE /home/mastodilu/Downloads/fiore_azzurro_1024x768.jpg
			WRITE /home/mastodilu/Downloads/fiore_azzurro_1024x768.jpg
			WRITE /home/mastodilu/Downloads/fiore_azzurro_1024x768.jpg
			CHMOD /home/mastodilu/Downloads/fiore_azzurro_1024x768.jpg

		se un file esiste gia' non viene lanciato CREATE, ma solo WRITE n volte seguito da CHMOD
	*/

}
