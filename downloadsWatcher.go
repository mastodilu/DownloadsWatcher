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

// CreteFolder creates a folder if it does not already exist
func CreteFolder(path string, perms uint32) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.FileMode(perms)); err != nil {
			log.Println(err)
		} else {
			fmt.Println("Created", path)
		}
	}
}

// selectDestFolder returns the right folder path based on the file extension passed as input
func selectDestFolder(sourceExtension string) string {
	sourceExtension = strings.TrimPrefix(sourceExtension, ".") // trim the . from the file extension
	for _, ext := range imageExtensions {
		if ext == sourceExtension {
			return filepath.Join(baseDst, "Pictures")
		}
	}
	for _, ext := range videoExtensions {
		if ext == sourceExtension {
			return filepath.Join(baseDst, "Videos")
		}
	}
	for _, ext := range audioExtensions {
		if ext == sourceExtension {
			return filepath.Join(baseDst, "Music")
		}
	}
	for _, ext := range documentExtensions {
		if ext == sourceExtension {
			return filepath.Join(baseDst, "Documents")
		}
	}
	for _, ext := range executableExtensions {
		if ext == sourceExtension {
			return filepath.Join(baseDst, "Downloads", "Eseguibili")
		}
	}
	for _, ext := range archivesExtensions {
		if ext == sourceExtension {
			return filepath.Join(baseDst, "Downloads", "Archivi")
		}
	}
	return filepath.Join(baseDst, "Downloads", "Altro")
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

	// create destination path string
	destPath := filepath.Join(selectDestFolder(fileextension), filename)

	// create new destination file into destination folder
	dstFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, os.FileMode(0666)) // TODO change baseDst based on file format
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

	fmt.Println("Moved", src, "into", destPath)
	return nil // OK
}

const (
	baseDst string = "/home/mastodilu"
)

var (
	imageExtensions      = []string{"jpeg", "jpg", "gif", "png", "bmp", "raw", "tiff", "psd", "cr2"}
	videoExtensions      = []string{"avi", "flv", "wmv", "mov", "mp4"}
	audioExtensions      = []string{"pcm", "wav", "aiff", "mp3", "aac", "ogg", "wma", "flac", "alac", "wma", "m3u"}
	documentExtensions   = []string{"pdf", "doc", "docx", "log", "odt", "rtf", "tex", "txt", "csv", "pps", "ppt", "pptx", "xml", "xls", "xlsx", "db", "dbf", "sql"}
	executableExtensions = []string{"apk", "app", "bat", "com", "exe", "jar", "wsf"}
	archivesExtensions   = []string{"7z", "deb", "tar.gz", "zip", "zipx"}
)

func main() {

	// create folders
	CreteFolder("/home/mastodilu/Downloads/Eseguibili", 0766)
	CreteFolder("/home/mastodilu/Downloads/Archivi", 0766)
	CreteFolder("/home/mastodilu/Downloads/Altro", 0766)

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

				// move file if chmod event occurs
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					if err := MoveFile(event.Name); err != nil {
						log.Println(err)
					}
				}

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
