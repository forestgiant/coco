package main

import (
	"strings"
	"testing"
)

func TestAddSpace(t *testing.T) {
	originalString := "ClientLoveProcess"
	spacedString := "Client Love Process"

	testString := addSpace(originalString)

	// check for spaces
	stringSpaces := strings.Split(testString, " ")
	if len(stringSpaces) < 2 {
		t.Fatalf("String not spaced")
	}

	if spacedString != testString {
		t.Fatalf("String does not match expected output")
	}
}
