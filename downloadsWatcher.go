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

	done := make(chan struct{})

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					// do something
				}

				fmt.Printf("%v %v\n",
					event.Op.String(), // operation name
					event.Name,        // file name
				)

			// watch for errors
			case err := <-watcher.Errors:
				done <- struct{}{} // terminate the program
				log.Fatalln(err)
			}
		}
	}()

	// watch for folder changes
	if err := watcher.Add("/home/mastodilu/Downloads"); err != nil {
		log.Fatalln(err)
	}

	<-done

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
