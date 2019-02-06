package main

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/roccoblues/tt/cmd"
)

const defaultFileName = ".tt.json"

func main() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}

	cfg := &cmd.Config{
		DefaultPath: filepath.Join(home, defaultFileName),
		RoundTo:     15,
		TimeFormat:  "15:04",
		DateFormat:  "02.01.2006",
	}

	c := cmd.New(cfg)
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
