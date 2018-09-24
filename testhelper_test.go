package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func newFile(t *testing.T, bytes []byte) string {
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpFile.Write(bytes); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}

func nonExistingFile(t *testing.T) string {
	path := newFile(t, []byte{})
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFile(t *testing.T, path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func newTime(t *testing.T, str string) time.Time {
	tm, err := time.ParseInLocation(dateTimeFormat, str, time.Now().Location())
	if err != nil {
		t.Fatal(err)
	}
	return tm
}
