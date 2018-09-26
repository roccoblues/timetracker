package main

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
)

const timeFormat = "15:04"
const dateFormat = "02.01.2006"

const defaultFileName = ".tt.json"
const defaultRoundToMinutes = 15

var dateTimeFormat = fmt.Sprintf("%s %s", dateFormat, timeFormat)

// flags
var path string
var roundToMinutes int

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}
	defaultPath := filepath.Join(home, defaultFileName)

	rootCmd.PersistentFlags().StringVarP(&path, "file", "f", defaultPath, "full path to data file")
	rootCmd.PersistentFlags().IntVarP(&roundToMinutes, "round-to", "r", defaultRoundToMinutes, "round times to minutes")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
