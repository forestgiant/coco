package main

import (
	"testing"
)

func TestAddSpace(t *testing.T) {
	originalString := "ClientLoveProcess"
	spacedString := "Client Love Process"

	testString := addSpace(originalString)

	if spacedString != testString {
		t.Error("String not spaced")
	}
}
