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

// FileName returns the filename of a given path or error
func FileName(src string) (string, error) {
	lastIndex := strings.LastIndexAny(src, "/")
	if lastIndex == -1 {
		return "", fmt.Errorf("last index of character '/' is -1 in source file")
	}
	filename := src[lastIndex:len(src)]
	return filename, nil
}

// FileExtension returns the file extension of a given filename
func FileExtension(filename string) string {
	lastIndex := strings.LastIndexAny(filename, ".")
	if lastIndex == -1 {
		return ""
	}
	ext := filename[lastIndex:len(filename)]
	return ext
}

// MoveFile moves src file into a destination folder
func MoveFile(src string) error {
	// open file to copy
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// obtain file name
	var filename string
	if filename, err = FileName(src); err != nil {
		return err
	}

	// obtain file extension
	fileextension := FileExtension(filename)

	fmt.Println("filename", filename, "extension", fileextension)

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

var (
	imageExtensions      = []string{"jpeg", "jpg", "gif", "png", "bmp", "raw", "tiff", "psd", "cr2"}
	videoExtensions      = []string{"avi", "flv", "wmv", "mov", "mp4"}
	audioExtensions      = []string{"pcm", "wav", "aiff", "mp3", "aac", "ogg", "wma", "flac", "alac", "wma", "m3u"}
	documentsExtensions  = []string{"doc", "docx", "log", "odt", "rtf", "tex", "txt", "csv", "pps", "ppt", "pptx", "xml", "xls", "xlsx", "db", "dbf", "sql"}
	executableExtensions = []string{"apk", "app", "bat", "com", "exe", "jar", "wsf"}
	// webFilesExtensions = []string{"asp", "aspx", "cer", "css", "htm", "html", "js", "jsp", "php", "xhtml"}
	archivesExtensions = []string{"7z", "deb", "tar.gz", "zip", "zipx"}
)

func main() {

	// create folders
	if err := os.Mkdir("/home/mastodilu/Downloads/Eseguibili", os.FileMode(0766)); err != nil {
		log.Println(err)
	}
	if err := os.Mkdir("/home/mastodilu/Downloads/Archivi", os.FileMode(0766)); err != nil {
		log.Println(err)
	}
	if err := os.Mkdir("/home/mastodilu/Downloads/Altro", os.FileMode(0766)); err != nil {
		log.Println(err)
	}

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
