package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "tt is a command line time tracker",
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := loadTimeSheet(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if month > 0 {
			ts.PrintMonth(time.Month(month), os.Stdout)
		} else {
			ts.Print(os.Stdout)
		}
	},
}
