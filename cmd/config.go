package cmd

import (
	"fmt"
	"time"
)

// Config is used to configure the commands.
type Config struct {
	TimeFormat  string
	DateFormat  string
	DefaultPath string
	Month       int
	RoundTo     int
	path        string
}

// DateTimeFormat returns the combined data and time format.
func (c *Config) DateTimeFormat() string {
	return fmt.Sprintf("%s %s", c.DateFormat, c.TimeFormat)
}

// RoundDuration return the RoundTo integer to a time.Duration.
func (c *Config) RoundDuration() time.Duration {
	return time.Duration(c.RoundTo) * time.Minute
}
