package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "Command line time tracking tool",
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := loadTimeSheet(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		ts.Print(os.Stdout, time.Duration(roundToMinutes)*time.Minute)
	},
}
