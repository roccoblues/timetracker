package timesheet

import (
	"fmt"
	"io"
	"time"
)

func print(times []time.Time, roundTo time.Duration, dateFormat, timeFormat string, out io.Writer) {
	// group times by day
	var days [][]time.Time
	var dayTimes []time.Time
	var prev time.Time

	for _, t := range times {
		if prev.IsZero() {
			prev = t
		}
		if !sameDate(prev, t) {
			days = append(days, dayTimes)
			dayTimes = []time.Time{}
		}
		dayTimes = append(dayTimes, t.Round(roundTo))
		prev = t
	}
	if len(dayTimes) > 0 {
		days = append(days, dayTimes)
	}

	if len(days) == 0 {
		return
	}

	var week int
	var totalHours time.Duration
	for _, times := range days {
		// output newline after each week
		_, w := times[0].ISOWeek()
		if week > 0 && week != w {
			fmt.Fprintln(out, "")
		}
		week = w

		// calculate hours per day
		var hours time.Duration
		var start time.Time
		for i, t := range times {
			if i%2 == 0 {
				start = t
			} else {
				hours += t.Sub(start)
			}
		}

		// output date and hours (ie. "01.09.2018 8.50")
		fmt.Fprintf(out, "%s  %.2f ", times[0].Format(dateFormat), hours.Hours())

		// output individual intervals (ie. "10:00-12:30 13:00-16:30")
		for i, t := range times {
			if i%2 == 0 {
				fmt.Fprintf(out, " %s-", t.Format(timeFormat))
			} else {
				fmt.Fprintf(out, "%s", t.Format(timeFormat))
			}
		}

		totalHours += hours

		fmt.Fprintln(out, "")
	}

	fmt.Fprintf(out, "\nTotal: %.2f\n", totalHours.Hours())
}
