package main

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func updateLinks(file []byte, folder os.FileInfo) []byte {
	re := regexp.MustCompile(`\([a-zA-Z]+.md\)`)

	links := re.FindAll(file, -1)

	var newFile []byte
	newFile = file

	removeDuplicates(&links)

	for _, link := range links {
		cleanedLink := strings.Split(string(link), "(")
		cleanedLink = strings.Split(cleanedLink[1], ")")
		fullLink := "/post/" + folder.Name() + "/" + cleanedLink[0]

		re = regexp.MustCompile(string(link))
		newFile = re.ReplaceAll(newFile, []byte(fullLink))
	}

	re = regexp.MustCompile(".md")

	return re.ReplaceAll(newFile, []byte(""))
}

func isReadme(file string) bool {
	if file == "README.md" || file == "readme.md" {
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

func removeDuplicates(xs *[][]byte) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[string(x)] {
			found[string(x)] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func generateHeader(title, category, date string) string {
	return "+++\ndate = \"" + date + "\"\ntitle = \"" + title + " - Table of Contents\"\ncategories = [\"" + category + "\"]\n\n+++"
}
