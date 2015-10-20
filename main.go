package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	contentDirectory     = flag.String("cf", "", "The content directory")
	hugoContentDirectory = flag.String("hf", "", "Hugo directory")
)

func main() {
	flag.Parse()

	if *contentDirectory == "" {
		fmt.Println("Please provide a content directory")
		os.Exit(1)
	}

	if *hugoContentDirectory == "" {
		fmt.Println("Please provide a hugo content directory")
		os.Exit(1)
	}

	// Get the content direcotry
	// Get all of the folders inside of it
	folders, err := ioutil.ReadDir(*contentDirectory)
	if err != nil {
		fmt.Println(err)
	}
	// Loop through each folder
	for _, folder := range folders {
		if folder.IsDir() {
			// Open the indexed folder and get all the files inside of it
			fullFolderPath := filepath.Join(*contentDirectory, folder.Name())
			hugoContentFolderPath := filepath.Join(*hugoContentDirectory, folder.Name())

			// create the hugo folder if it doesn't exist
			err := os.MkdirAll(hugoContentFolderPath, 0777)
			if err != nil {
				fmt.Println(err)
			}

			// get all the files inside the indexed folder
			files, err := ioutil.ReadDir(fullFolderPath)
			if err != nil {
				fmt.Println(err)
			}

			for _, file := range files {
				if !file.IsDir() {
					header := "+++ \n date = \"2015-10-15\" \n title = \"" + file.Name() + "\" \n+++"

					indexedFilePath := filepath.Join(fullFolderPath, file.Name())

					readFile, err := ioutil.ReadFile(indexedFilePath)
					if err != nil {
						fmt.Println(err)
					}

					concatFile := header + "\n \n" + string(readFile)

					tmpFilePath := filepath.Join(folder.Name(), file.Name())
					fullFilePath := filepath.Join(*hugoContentDirectory, tmpFilePath)

					// create the file
					_, err = os.Create(fullFilePath)
					if err != nil {
						fmt.Println(err)
					}

					// write the file
					err = ioutil.WriteFile(fullFilePath, []byte(concatFile), 0644)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}

	fmt.Println("All Done =]")
}
