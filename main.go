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
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(home, fileName)

	ts := &timeSheet{}

	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		bytes, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(bytes, &ts); err != nil {
			return err
		}
	}

	modified := false

	if len(args) > 1 {
		cmd := args[1]
		switch cmd {
		case "start":
			if err := ts.Start(time.Now()); err != nil {
				return err
			}
			modified = true
		case "stop":
			if err := ts.End(time.Now()); err != nil {
				return err
			}
			modified = true
		default:
			return fmt.Errorf("unknown command: %s", cmd)
		}
	}

	if modified {
		bytes, err := json.MarshalIndent(ts, "", "  ")
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(fullPath, bytes, 0644); err != nil {
			return err
		}
	}

	writeDays(ts.Days(), os.Stdout)

	return nil
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
