// This script allows to track the downloads folder for changes.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// MoveFile moves src file into a destination folder
func MoveFile(src string) error {
	// open file to copy
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	lastIndex := strings.LastIndexAny(src, "/")
	if lastIndex == -1 {
		return fmt.Errorf("last index of character '/' is -1 in source file")
	}
	filename := src[lastIndex:len(src)]

	// create new destination file into destination folder
	dstFile, err := os.OpenFile(filepath.Join(baseDst, filename), os.O_CREATE|os.O_WRONLY|os.O_EXCL, os.FileMode(0666)) // TODO change baseDst based on file format
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// write into destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// delete original file
	err = os.Remove(src)
	if err != nil {
		return err
	}

	fmt.Println("Moved", src)
	return nil // OK
}

const (
	baseDst string = "/home/mastodilu/PROVA"
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
					err := MoveFile(event.Name)
					if err != nil {
						log.Fatalln(err)
					}
				}

				// fmt.Printf("%v %v\n",
				// 	event.Op.String(), // operation name
				// 	event.Name,        // file name
				// )

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
