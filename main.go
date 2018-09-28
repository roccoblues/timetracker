package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

const timeFormat = "15:04"
const dateFormat = "02.01.2006"
const roundTo = 15 * time.Minute
const defaultFileName = ".tt.json"

var dateTimeFormat = fmt.Sprintf("%s %s", dateFormat, timeFormat)

// flags
var file string

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}
	defaultPath := filepath.Join(home, defaultFileName)

	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", defaultPath, "path to data `FILE`")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
