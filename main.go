package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/roccoblues/tt/timesheet"
)

const defaultFileName = ".tt.json"

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [start|stop] [time]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flagFile := flag.String("file", filepath.Join(home, defaultFileName), "path to data file")
	flagMonth := flag.Int("month", 0, "output month (default current)")
	flagDateFormat := flag.String("date-format", "02.01.2006", "parse and write dates with format")
	flagTimeFormat := flag.String("time-format", "15:04", "parse and write times with format")
	flagRoundTo := flag.Int("round-to", 15, "round to minutes")
	flag.Parse()

	var month time.Month
	if *flagMonth == 0 {
		month = time.Now().Month()
	} else {
		month = time.Month(*flagMonth)
	}

	timeArg := time.Now()
	if len(flag.Arg(1)) != 0 {
		timeArg, err = parseTime(flag.Arg(1), *flagDateFormat, *flagTimeFormat)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	sheet := loadSheet(*flagFile, *flagDateFormat, *flagTimeFormat)

	if len(flag.Args()) == 0 {
		sheet.PrintMonth(month, time.Duration(*flagRoundTo)*time.Minute, os.Stdout)
		os.Exit(0)
	}

	switch flag.Arg(0) {
	default:
		fmt.Fprintf(os.Stderr, "%s: unknown command '%s'\n", os.Args[0], flag.Arg(0))
		os.Exit(1)
	case "start":
		if err := sheet.Start(timeArg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "stop":
		if err := sheet.End(timeArg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	saveSheet(sheet, *flagFile)
	sheet.PrintMonth(month, time.Duration(*flagRoundTo)*time.Minute, os.Stdout)
}

func loadSheet(path string, dateFormat, timeFormat string) *timesheet.Sheet {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	s, err := timesheet.Load(f, dateFormat, timeFormat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return s
}

func saveSheet(s *timesheet.Sheet, path string) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	if err := s.Save(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseTime(value string, dateFormat, timeFormat string) (time.Time, error) {
	dateTimeFormat := fmt.Sprintf("%s %s", dateFormat, timeFormat)
	dateTime := fmt.Sprintf("%s %s", time.Now().Format(dateFormat), value)

	return time.ParseInLocation(dateTimeFormat, dateTime, time.Now().Location())
}
