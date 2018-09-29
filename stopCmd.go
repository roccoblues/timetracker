package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [time]",
	Short: "Stop the current timetracking interval",
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := loadTimeSheet(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		endTime := time.Now()
		if len(args) > 0 {
			endTime, err = parseTime(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		if err := ts.End(endTime); err != nil {
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
