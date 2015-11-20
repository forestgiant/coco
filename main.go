package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Coco has the folder paths needed
// to generate the necessary Hugo files
type Coco struct {
	// The process directory
	ProcessDir string
	// The hugo directory
	HugoDir string
	// The indexed folder inside of the process dir
	ProcessIndexDir string
	// The hugo content
	HugoContentDir string
}

var (
	pushPtr = flag.Bool("push", false, "push process site to github")
)

func main() {
	// Parse the passed in flags and args
	flag.Parse()

	// MD Hugo version
	// version := "0.1.1"

	// Validate the correct args were passed in
	if len(flag.Args()) < 2 {
		fmt.Println("Please pass in the required arguments")
		os.Exit(1)
	}

	// Get the os args
	processDirectory := flag.Args()[0]
	hugoDirectory := flag.Args()[1]

	// Instantiate a Coco
	coco := &Coco{
		ProcessDir: processDirectory,
		HugoDir:    hugoDirectory,
	}

	// Read the process directory
	folders, err := ioutil.ReadDir(processDirectory)
	if err != nil {
		fmt.Println("Could not read process directory:", err)
	}

	// Loop through each folder
	for _, folder := range folders {
		// check to make sure folder is a directory
		if folder.IsDir() {
			coco.run(folder)
		}
	}

	// Push process website to GitHub
	if *pushPtr == true {
		coco.push()
	}

	fmt.Println("All done =]")
}

func (c *Coco) run(folder os.FileInfo) {

	// Set up the folder paths
	c.ProcessIndexDir = filepath.Join(c.ProcessDir, folder.Name())
	c.HugoContentDir = filepath.Join(c.HugoDir, folder.Name())

	// Create the folder in Hugo if it does not exist
	if err := os.MkdirAll(c.HugoContentDir, 0777); err != nil {
		fmt.Println("Error creating folder in hugo content:", err)
	}

	// Read the processIndexDir
	files, err := ioutil.ReadDir(c.ProcessIndexDir)
	if err != nil {
		fmt.Println("Error reading process index directory:", err)
	}

	// Loop through each folder/file in processIndexDir
	for _, file := range files {
		// Make sure its a file
		if !file.IsDir() {
			c.createFile(file, folder)
		}
	}
}

func (c *Coco) createFile(file, folder os.FileInfo) {

	// Check if it's the README.md file
	if !isReadme(file.Name()) {
		// Sanitize title
		sanitizedTitle := sanitizeTitle(file.Name())

		// Generate the header
		header := generateHeader(sanitizedTitle, folder.Name(), string(file.ModTime().Format(time.RFC3339)))

		// Generate the file content
		c.generateFileContent(header, file, folder)
	} else {
		// Generate the header
		header := generateHeader(strings.Title(folder.Name())+"- Table of Contents", folder.Name(), string(time.Now().Format(time.RFC3339)))

		// Generate the file content
		c.generateFileContent(header, file, folder)
	}
}

func (c *Coco) generateFileContent(header string, file, folder os.FileInfo) {
	// Get the filepath and read it
	processFilePath := filepath.Join(c.ProcessIndexDir, file.Name())
	readFile, err := ioutil.ReadFile(processFilePath)
	if err != nil {
		fmt.Println("Error reading process file:", err)
	}

	// Update the links to work with site
	updatedFile := updateLinks(readFile, folder)

	// Add the header to the file
	cleanedFile := header + "\n\n" + string(updatedFile)

	tmpFilePath := filepath.Join(folder.Name(), file.Name())
	hugoFilePath := filepath.Join(c.HugoDir, tmpFilePath)

	// create the file
	_, err = os.Create(hugoFilePath)
	if err != nil {
		fmt.Println(err)
	}

	// write the file
	err = ioutil.WriteFile(hugoFilePath, []byte(cleanedFile), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}

func (c *Coco) push() {
	// The main Hugo src
	mainHugoFolderPath := c.HugoDir + "/../.."

	// Switch folder path to mainHugoFolderPath
	if err := os.Chdir(mainHugoFolderPath); err != nil {
		fmt.Println("Error switching directories", err)
	}

	// Run git add .
	cmd := exec.Command("git", "add", ".")
	if err := cmd.Run(); err != nil {
		fmt.Println("Could not add git", err)
		log.Fatal(err)
	}

	// Run git commit
	out, err := exec.Command("git", "commit", "-mupdate site").Output()
	if err != nil {
		fmt.Println(string(out))
		log.Fatal(err)
	}

	// Run git push
	out, err = exec.Command("git", "push").Output()
	if err != nil {
		fmt.Println(string(out))
		log.Fatal(err)
	}

}
