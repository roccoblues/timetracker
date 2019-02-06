package timesheet

import (
	"io/ioutil"
	"testing"
)

func readFile(t *testing.T, path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
