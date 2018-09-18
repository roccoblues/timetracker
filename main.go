package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

const fileName = ".tt.json"
const roundTo = 15 * time.Minute

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find home directory: %v\n", err)
		os.Exit(1)
	}

	storage := newFileRepo(filepath.Join(home, fileName))

	tracker, err := newTracker(storage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create tracker: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "start":
			if err := tracker.Start(time.Now()); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to add start time: %v\n", err)
				os.Exit(1)
			}
		case "stop":
			if err := tracker.End(time.Now()); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to add end time: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
			os.Exit(1)
		}
	}

	fmt.Print(formatDays(tracker.Days()))
}

func formatDays(days []*day) string {
	var b bytes.Buffer
	var week int
	for _, day := range days {
		_, w := day.Date.ISOWeek()
		if week > 0 && week != w {
			b.WriteString("\n")
		}
		week = w

		b.WriteString(fmt.Sprintf("%s  %.2f  ", day.Date.Format("02.01.2006"), day.Time().Round(roundTo).Hours()))

		for i, e := range day.Entries {
			start := e.Start.Round(roundTo).Format("15:04")
			if e.End.IsZero() {
				b.WriteString(fmt.Sprintf("%s-", start))
			} else {
				end := e.End.Round(roundTo).Format("15:04")
				b.WriteString(fmt.Sprintf("%s-%s ", start, end))
			}

			if i == len(day.Entries)-1 {
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}
