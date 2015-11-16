package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
)

var (
	pushPtr = flag.Bool("push", false, "to push or not to push")
)

func main() {

	flag.Parse()

	contentDirectory := flag.Args()[0]
	hugoContentDirectory := flag.Args()[1]

	// Create wait group for run function
	var wg sync.WaitGroup

	// Create wait group for createFile function
	var fwg sync.WaitGroup

	// Quick valildation
	if contentDirectory == "" {
		fmt.Println("Please provide a content directory")
		os.Exit(2)
	}

	if hugoContentDirectory == "" {
		fmt.Println("Please provide a hugo content directory")
		os.Exit(2)
	}

	// Get the content direcotry
	// Get all of the folders inside of it
	folders, err := ioutil.ReadDir(contentDirectory)
	if err != nil {
		fmt.Println(err)
	}

	// Loop through each folder
	for _, folder := range folders {
		if folder.IsDir() {
			wg.Add(1)
			go run(folder, contentDirectory, hugoContentDirectory, &wg, &fwg)
		}
	}

	wg.Wait()

	if *pushPtr == true {
		mainHugoFolderPath := hugoContentDirectory + "/.."

		if err := os.Chdir(mainHugoFolderPath); err != nil {
			log.Print("Error switching directories", err)
		}

		cmd := exec.Command("git", "add", ".")
		if err := cmd.Run(); err != nil {
			fmt.Println("Could not add git", err)
			log.Fatal(err)
		}

		out, err := exec.Command("git", "commit", "-m=\"update site\"").Output()
		if err != nil {
			fmt.Println(string(out))
			log.Fatal(err)
		}

		cmd = exec.Command("git", "push")
		if err := cmd.Run(); err != nil {
			fmt.Println("could not push", err)
			log.Fatal(err)
		}
	}

	fmt.Println("All done =]")
}

func run(folder os.FileInfo, contentDirectory, hugoContentDirectory string, wg, fwg *sync.WaitGroup) {
	defer wg.Done()

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
			fwg.Add(1)
			go createFile(file, folder, indexedFolderPath, hugoContentDirectory, fwg)
		}
	}

	fwg.Wait()
}

func createFile(file, folder os.FileInfo, indexedFolderPath, hugoContentDirectory string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Create the sanitized title
	// from the file name
	if !isReadme(file) {
		sanitizedTitle := sanitizeTitle(file.Name())

		// Header added at the top of
		// every Hugo page
		header := "+++\ndate = \"" + string(file.ModTime().Format(time.RFC3339)) + "\"\ntitle = \"" + sanitizedTitle + "\"\ncategories = [\"" + folder.Name() + "\"]\n\n+++"

		// Get the file path and read it
		indexedFilePath := filepath.Join(indexedFolderPath, file.Name())
		readFile, err := ioutil.ReadFile(indexedFilePath)
		if err != nil {
			fmt.Println(err)
		}

		re := regexp.MustCompile(`\([a-zA-Z]+.md\)`)

		links := re.FindAll(readFile, -1)

		var newLink []byte

		for _, link := range links {
			cleanedLink := strings.Split(string(link), "(")
			fullLink := "/post/" + folder.Name() + "/" + cleanedLink[1]

			re = regexp.MustCompile(string(link))
			newLink = re.ReplaceAll(readFile, []byte(fullLink))
		}

		re = regexp.MustCompile(".md")

		newFile := re.ReplaceAll(newLink, []byte(""))

		// Add the header to the file
		concatFile := header + "\n \n" + string(newFile)

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
	} else {

		// Header added at the top of
		// every Hugo page
		header := "+++\ndate = \"" + string(time.Now().Format(time.RFC3339)) + "\"\ntitle = \"" + strings.Title(folder.Name()) + " - Table of Contents\"\ncategories = [\"" + folder.Name() + "\"]\n\n+++"

		// Get the file path and read it
		indexedFilePath := filepath.Join(indexedFolderPath, file.Name())
		readFile, err := ioutil.ReadFile(indexedFilePath)
		if err != nil {
			fmt.Println(err)
		}

		re := regexp.MustCompile(`\([a-zA-Z]+.md\)`)

		links := re.FindAll(readFile, -1)

		var newLink []byte

		for _, link := range links {
			cleanedLink := strings.Split(string(link), "(")
			fullLink := "/post/" + folder.Name() + "/" + cleanedLink[1]

			re = regexp.MustCompile(string(link))
			newLink = re.ReplaceAll(readFile, []byte(fullLink))
		}

		re = regexp.MustCompile(".md")

		newFile := re.ReplaceAll(newLink, []byte(""))

		// Add the header to the file
		concatFile := header + "\n \n" + string(newFile)

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

func isReadme(file os.FileInfo) bool {
	if file.Name() == "README.md" || file.Name() == "readme.md" {
		return true
	}

	return false
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
