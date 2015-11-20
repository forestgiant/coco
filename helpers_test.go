package main

import (
	"strings"
	"testing"
)

func TestSanitizeTitle(t *testing.T) {
	fileName := "TestFile.md"
	sanitizedTitle := "Test File"

	testString := sanitizeTitle(fileName)

	if testString != sanitizedTitle {
		t.Error("String not sanitized")
	}
}

func TestAddSpace(t *testing.T) {
	originalString := "ClientLoveProcess"
	spacedString := "Client Love Process"

	testString := addSpace(originalString)

	// check for spaces
	stringSpaces := strings.Split(testString, " ")
	if len(stringSpaces) < 2 {
		t.Error("String not spaced")
	}

	if spacedString != testString {
		t.Error("String does not match expected output")
	}
}
