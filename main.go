package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

const fileName = ".tt.json"
const roundTo = 15 * time.Minute
const timeFormat = "15:04"
const dateTimeFormat = "2006-01-02 15:04"

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find home directory: %v\n", err)
		os.Exit(1)
	}
	fullPath := filepath.Join(home, fileName)

	ts := &timeSheet{}

	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		bytes, err := ioutil.ReadFile(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read file '%s': %v\n", fullPath, err)
			os.Exit(1)
		}
		if err := json.Unmarshal(bytes, &ts); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode json: %v\n", err)
			os.Exit(1)
		}
	}

	modified := false

	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "start":
			if err := ts.Start(time.Now()); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to add start time: %v\n", err)
				os.Exit(1)
			}
			modified = true
		case "stop":
			if err := ts.End(time.Now()); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to add end time: %v\n", err)
				os.Exit(1)
			}
			modified = true
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
			os.Exit(1)
		}
	}

	if modified {
		bytes, err := json.MarshalIndent(ts, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to encode json: %v\n", err)
			os.Exit(1)
		}

		if err := ioutil.WriteFile(fullPath, bytes, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write file '%s': %v\n", fullPath, err)
			os.Exit(1)
		}
	}

	writeDays(ts.Days(), os.Stdout)
}

func writeDays(days []*day, output io.Writer) {
	var week int
	for _, day := range days {
		_, w := day.Date.ISOWeek()
		if week > 0 && week != w {
			fmt.Fprintln(output, "")
		}
		week = w

		fmt.Fprintf(output, "%s  %.2f  ", day.Date.Format("02.01.2006"), day.Time().Round(roundTo).Hours())

		for i, e := range day.Entries {
			start := e.Start.Round(roundTo).Format(timeFormat)
			if e.End.IsZero() {
				fmt.Fprintf(output, "%s-", start)
			} else {
				end := e.End.Round(roundTo).Format(timeFormat)
				fmt.Fprintf(output, "%s-%s ", start, end)
			}

			if i == len(day.Entries)-1 {
				fmt.Fprintln(output, "")
			}
		}
	}
}
