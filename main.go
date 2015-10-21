package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

func main() {

	contentDirectory := os.Args[1]
	hugoContentDirectory := os.Args[2]

	// Channel for when function is done
	exit := make(chan bool)

	// Quick valildation
	if contentDirectory == "" {
		fmt.Println("Please provide a content directory")
		os.Exit(1)
	}

	if hugoContentDirectory == "" {
		fmt.Println("Please provide a hugo content directory")
		os.Exit(1)
	}

	// Get the content direcotry
	// Get all of the folders inside of it
	folders, err := ioutil.ReadDir(contentDirectory)
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		// Loop through each folder
		for _, folder := range folders {
			if folder.IsDir() {
				// Open the indexed folder and get all the files inside of it
				indexedFolderPath := filepath.Join(contentDirectory, folder.Name())
				hugoContentFolderPath := filepath.Join(hugoContentDirectory, folder.Name())

				// create the hugo folder if it doesn't exist
				err := os.MkdirAll(hugoContentFolderPath, 0777)
				if err != nil {
					fmt.Println(err)
				}

				// get all the files inside the indexed folder
				files, err := ioutil.ReadDir(indexedFolderPath)
				if err != nil {
					fmt.Println(err)
				}

				for _, file := range files {
					if !file.IsDir() {
						// Create the sanitized title
						// from the file name
						sanitizedTitle := sanitizeTitle(file.Name())

						// Header added at the top of
						// every Hugo page
						header := "+++ \n date = \"" + string(file.ModTime().Format(time.UnixDate)) + "\" \n title = \"" + sanitizedTitle + "\" \n+++"

						// Get the file path and read it
						indexedFilePath := filepath.Join(indexedFolderPath, file.Name())
						readFile, err := ioutil.ReadFile(indexedFilePath)
						if err != nil {
							fmt.Println(err)
						}

						// Add the header to the file
						concatFile := header + "\n \n" + string(readFile)

						tmpFilePath := filepath.Join(folder.Name(), file.Name())
						hugoFilePath := filepath.Join(hugoContentDirectory, tmpFilePath)

						// create the file
						_, err = os.Create(hugoFilePath)
						if err != nil {
							fmt.Println(err)
						}

						// write the file
						err = ioutil.WriteFile(hugoFilePath, []byte(concatFile), 0644)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}

		exit <- true
	}()

	<-exit

	fmt.Println("All done =]")
}

func sanitizeTitle(s string) string {
	splitName := strings.Split(s, ".")
	fileName := splitName[0]
	return addSpace(fileName)
}

func addSpace(s string) string {
	buf := &bytes.Buffer{}
	for i, rune := range s {
		if unicode.IsUpper(rune) && i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteRune(rune)
	}
	return buf.String()
}
