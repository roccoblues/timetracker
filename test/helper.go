package test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
)

func NewFile(t *testing.T, data []byte) string {
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpFile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpFile.Name()
}

func NonExistingFile(t *testing.T) string {
	tmpFile := NewFile(t, []byte{})
	if err := os.Remove(tmpFile); err != nil {
		t.Fatal(err)
	}
	return tmpFile
}

func ReadFile(t *testing.T, path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func Time(t *testing.T, str string) time.Time {
	tm, err := time.ParseInLocation("2006-01-02 15:04", str, time.Now().Location())
	if err != nil {
		t.Fatal(err)
	}
	return tm
}

var InvalidJSON = "{"

var EmptyJSON = "{}"

var EmptyDayJSON = heredoc.Doc(`{
  "2018-09-01": []
}`)

var OneDayOnlyStartJSON = heredoc.Doc(`{
  "2018-09-01": [
    "10:00"
  ]
}`)

var OneDayStartEndJSON = heredoc.Doc(`{
  "2018-09-01": [
    "10:00",
    "12:00"
  ]
}`)

var OneDayStartEndStartJSON = heredoc.Doc(`{
  "2018-09-01": [
    "10:00",
    "12:00",
    "13:00"
  ]
}`)

var MultipleDaysJSON = heredoc.Doc(`{
  "2018-09-01": [
    "10:00",
    "12:00"
  ],
  "2018-09-02": [
    "08:00"
  ]
}`)
var MultipleDaysEndJSON = heredoc.Doc(`{
  "2018-09-01": [
    "10:00",
    "12:00"
  ],
  "2018-09-02": [
    "08:00",
    "16:00"
  ]
}`)

var InvalidDateJSON = heredoc.Doc(`{
  "2018-09-00": [
    "10:00"
  ]
}`)

var InvalidTimeJSON = heredoc.Doc(`{
  "2018-09-01": [
    "25:00"
  ]
}`)
