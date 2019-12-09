package timesheet

import (
	"fmt"
	"io"
	"time"
)

func print(times []time.Time, roundTo time.Duration, dateFormat, timeFormat string, out io.Writer) {
	for i, t := range times {
		times[i] = t.Round(roundTo)
	}

	days := groupTimesByDay(times)

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

		hours := calculateHours(times)

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

func groupTimesByDay(times []time.Time) [][]time.Time {
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
		dayTimes = append(dayTimes, t)
		prev = t
	}
	if len(dayTimes) > 0 {
		days = append(days, dayTimes)
	}

	return days
}

func calculateHours(times []time.Time) time.Duration {
	var hours time.Duration

	var start time.Time
	for i, t := range times {
		if i%2 == 0 {
			start = t
		} else {
			hours += t.Sub(start)
		}
	}

	return hours
}
