package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

const timeFormat = "15:04"
const dateFormat = "02.01.2006"
const defaultCmd = "print"

const defaultFileName = ".tt.json"
const defaultRoundToMinutes = 15

var dateTimeFormat = fmt.Sprintf("%s %s", dateFormat, timeFormat)

type config struct {
	path    string
	roundTo time.Duration
	out     io.Writer
}

type command interface {
	Name() string
	Description() string
	Run(*config) error
	Default() bool
}

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}
	defaultPath := filepath.Join(home, defaultFileName)

	// Redefining Usage() customizes the output of `tt -h`
	flag.Usage = printUsage

	var (
		path           = flag.String("file", defaultPath, "full path to data file")
		roundToMinutes = flag.Int("round-to", defaultRoundToMinutes, "round times to minutes")
	)
	flag.Parse()

	config := config{
		path:    *path,
		roundTo: time.Duration(*roundToMinutes) * time.Minute,
		out:     os.Stdout,
	}

	var cmdName string
	runDefaultCmd := false
	if flag.NArg() == 0 {
		runDefaultCmd = true
	} else {
		cmdName = flag.Args()[0]
	}

	for _, cmd := range commandList() {
		if runDefaultCmd && cmd.Default() {
			if err := cmd.Run(&config); err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
		if cmdName == cmd.Name() {
			if err := cmd.Run(&config); err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	// unknown command, print usage help and exit
	printUsage()
	os.Exit(1)
}

var usageText = `Command line time tracking tool

Usage:
  tt [command]

Available Commands:
`

func printUsage() {
	fmt.Fprintln(flag.CommandLine.Output(), usageText)
	for _, cmd := range commandList() {
		fmt.Fprintf(flag.CommandLine.Output(), "  %s\t %s", cmd.Name(), cmd.Description())
		if cmd.Default() {
			fmt.Fprintf(flag.CommandLine.Output(), " (default command)")
		}
		fmt.Fprintln(flag.CommandLine.Output(), "")
	}
	fmt.Fprintln(flag.CommandLine.Output(), "\nFlags:")
	flag.PrintDefaults()
}

func commandList() []command {
	return []command{
		&startCmd{},
		&stopCmd{},
		&printCmd{},
	}
}
