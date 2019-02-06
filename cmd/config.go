package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/roccoblues/tt/timesheet"
)

type Config struct {
	TimeFormat  string
	DateFormat  string
	DefaultPath string
	Month       int
	RoundTo     int
	path        string
}

func (c *Config) DateTimeFormat() string {
	return fmt.Sprintf("%s %s", c.DateFormat, c.TimeFormat)
}

func (c *Config) RoundDuration() time.Duration {
	return time.Duration(c.RoundTo) * time.Minute
}

func (c *Config) parseTime(value string) (time.Time, error) {
	dateTime := fmt.Sprintf("%s %s", time.Now().Format(c.DateFormat), value)
	return time.ParseInLocation(c.DateTimeFormat(), dateTime, time.Now().Location())
}

func (c *Config) loadSheet() *timesheet.Sheet {
	f, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	s, err := timesheet.Load(f, c.DateFormat, c.TimeFormat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return s
}

func (c *Config) saveSheet(s *timesheet.Sheet) {
	f, err := os.Create(c.path)
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
