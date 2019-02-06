package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/roccoblues/tt/timesheet"
	"github.com/spf13/cobra"
)

// New returns a cobra rootCmd configured with the given config.
func New(cfg *Config) *cobra.Command {
	rootCmd := newRootCommand(cfg)

	rootCmd.Flags().IntVarP(&cfg.Month, "month", "m", 0, "output `MONTH` (default current)")

	rootCmd.PersistentFlags().StringVarP(&cfg.path, "file", "f", "", "path to data `FILE` (default $HOME/.tt.json)")
	rootCmd.PersistentFlags().IntVarP(&cfg.RoundTo, "round-to", "r", cfg.RoundTo, "round to `MINUTES`")
	rootCmd.PersistentFlags().StringVarP(&cfg.TimeFormat, "time-format", "t", cfg.TimeFormat, "parse and write times with `FORMAT`")
	rootCmd.PersistentFlags().StringVarP(&cfg.DateFormat, "date-format", "d", cfg.DateFormat, "parse and write dates with `FORMAT`")

	if !rootCmd.Flags().Changed("month") {
		cfg.Month = int(time.Now().Month())
	}
	if !rootCmd.Flags().Changed("path") {
		cfg.path = cfg.DefaultPath
	}

	rootCmd.AddCommand(newStartCommand(cfg))
	rootCmd.AddCommand(newStopCmd(cfg))

	return rootCmd
}

func newRootCommand(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "tt",
		Short: "tt is a command line time tracker",
		Run: func(cmd *cobra.Command, args []string) {
			sheet := loadSheet(cfg)
			sheet.PrintMonth(time.Month(cfg.Month), cfg.RoundDuration(), os.Stdout)
		},
	}
}

func newStartCommand(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "start [time]",
		Short: "Start a new timetracking interval",
		Run: func(cmd *cobra.Command, args []string) {
			sheet := loadSheet(cfg)

			if err := sheet.Start(timeArg(args, cfg)); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			sheet.PrintMonth(time.Month(cfg.Month), cfg.RoundDuration(), os.Stdout)

			saveSheet(sheet, cfg)
		},
	}
}

func newStopCmd(cfg *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [time]",
		Short: "Stop the current timetracking interval",
		Run: func(cmd *cobra.Command, args []string) {
			sheet := loadSheet(cfg)

			if err := sheet.End(timeArg(args, cfg)); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			sheet.PrintMonth(time.Month(cfg.Month), cfg.RoundDuration(), os.Stdout)

			saveSheet(sheet, cfg)
		},
	}
}

func timeArg(args []string, cfg *Config) time.Time {
	t := time.Now()
	var err error

	if len(args) > 0 {
		t, err = parseTime(args[0], cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	return t
}

func parseTime(value string, cfg *Config) (time.Time, error) {
	dateTimeFormat := fmt.Sprintf("%s %s", cfg.DateFormat, cfg.TimeFormat)
	dateTime := fmt.Sprintf("%s %s", time.Now().Format(cfg.DateFormat), value)

	return time.ParseInLocation(dateTimeFormat, dateTime, time.Now().Location())
}

func loadSheet(cfg *Config) *timesheet.Sheet {
	f, err := os.OpenFile(cfg.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	s, err := timesheet.Load(f, cfg.DateFormat, cfg.TimeFormat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return s
}

func saveSheet(s *timesheet.Sheet, cfg *Config) {
	f, err := os.Create(cfg.path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	if err := s.Save(f); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
