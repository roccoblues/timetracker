package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func New(cfg *Config) *cobra.Command {
	rootCmd := newRootCommand(cfg)

	rootCmd.Flags().IntVarP(&cfg.Month, "month", "m", 0, "output `MONTH` (default current)")
	if !rootCmd.Flags().Changed("month") {
		cfg.Month = int(time.Now().Month())
	}

	rootCmd.PersistentFlags().StringVarP(&cfg.Path, "file", "f", cfg.Path, "path to data `FILE`")
	rootCmd.PersistentFlags().IntVarP(&cfg.RoundTo, "round-to", "r", cfg.RoundTo, "round to `MINUTES`")
	rootCmd.PersistentFlags().StringVarP(&cfg.TimeFormat, "time-format", "t", cfg.TimeFormat, "parse and write times with `FORMAT`")
	rootCmd.PersistentFlags().StringVarP(&cfg.DateFormat, "date-format", "d", cfg.DateFormat, "parse and write dates with `FORMAT`")

	rootCmd.AddCommand(newStartCommand(cfg))
	rootCmd.AddCommand(newStopCmd(cfg))

	return rootCmd
}

func newRootCommand(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "tt",
		Short: "tt is a command line time tracker",
		Run: func(cmd *cobra.Command, args []string) {
			s := cfg.loadSheet()
			s.PrintMonth(time.Month(cfg.Month), cfg.RoundDuration(), os.Stdout)
		},
	}

}

func newStartCommand(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "start [time]",
		Short: "Start a new timetracking interval",
		Run: func(cmd *cobra.Command, args []string) {
			sheet := cfg.loadSheet()

			startTime := time.Now()
			var err error

			if len(args) > 0 {
				startTime, err = cfg.parseTime(args[0])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			if err := sheet.Start(startTime); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			cfg.saveSheet(sheet)

			sheet.Print(cfg.RoundDuration(), os.Stdout)
		},
	}
}

func newStopCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [time]",
		Short: "Stop the current timetracking interval",
		Run: func(cmd *cobra.Command, args []string) {
			sheet := cfg.loadSheet()

			endTime := time.Now()
			var err error

			if len(args) > 0 {
				endTime, err = cfg.parseTime(args[0])
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			if err := sheet.End(endTime); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			cfg.saveSheet(sheet)

			sheet.Print(cfg.RoundDuration(), os.Stdout)
		},
	}
}
