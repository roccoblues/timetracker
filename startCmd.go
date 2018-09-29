package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [time]",
	Short: "Start a new timetracking interval",
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := loadTimeSheet(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		startTime := time.Now()
		if len(args) > 0 {
			startTime, err = parseTime(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		if err := ts.Start(startTime); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := ts.Save(file); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		ts.Print(os.Stdout)
	},
}
