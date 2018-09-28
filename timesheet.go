package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type timeSheet struct {
	times []time.Time
}

func loadTimeSheet(path string) (*timeSheet, error) {
	ts := timeSheet{}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(bytes, &ts); err != nil {
			return nil, err
		}
	}

	return &ts, nil
}

func (ts *timeSheet) Start(start time.Time) error {
	times := ts.TimesForDay(start)

	// an uneven number of times means we already started an interval
	if len(times)%2 != 0 {
		return fmt.Errorf("already started")
	}

	if len(times) > 0 {
		last := times[len(times)-1]
		if start.Before(last) {
			return fmt.Errorf("start time %s is ealier as last end time %s", start.Format(timeFormat), last.Format(timeFormat))
		}
	}

	ts.times = append(ts.times, start)

	return nil
}

func (ts *timeSheet) End(end time.Time) error {
	times := ts.TimesForDay(end)

	// an even number of times means we haven't started an interval
	if len(times)%2 == 0 {
		return fmt.Errorf("not started")
	}

	if len(times) > 0 {
		last := times[len(times)-1]
		if end.Before(last) {
			return fmt.Errorf("end time %s is earlier as last start time %s", end.Format(timeFormat), last.Format(timeFormat))
		}
	}

	ts.times = append(ts.times, end)

	return nil
}

func (ts *timeSheet) MarshalJSON() ([]byte, error) {
	data := map[string][]string{}

	for _, t := range ts.times {
		date := t.Format(dateFormat)
		if _, exists := data[date]; !exists {
			data[date] = []string{}
		}
		data[date] = append(data[date], t.Format(timeFormat))
	}

	return json.MarshalIndent(&data, "", "  ")
}

func (ts *timeSheet) UnmarshalJSON(bytes []byte) error {
	var decoded map[string][]string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return err
	}

	dates := make([]string, 0, len(decoded))
	for d := range decoded {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	loc := time.Now().Location()

	for _, d := range dates {
		for _, e := range decoded[d] {
			timeString := fmt.Sprintf("%s %s", d, e)
			tm, err := time.ParseInLocation(dateTimeFormat, timeString, loc)
			if err != nil {
				return err
			}
			ts.times = append(ts.times, tm)
		}
	}

	return nil
}

func (ts *timeSheet) Save(path string) error {
	bytes, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (ts *timeSheet) Print(out io.Writer) {
	// group times by day
	var days [][]time.Time
	var prev time.Time
	i := 0
	for _, t := range ts.times {
		if prev.IsZero() {
			prev = t
		}
		if !sameDate(prev, t) {
			i++
		}
		if i > len(days)-1 {
			days = append(days, []time.Time{})
		}
		days[i] = append(days[i], t.Round(roundTo))
		prev = t
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

func (ts *timeSheet) TimesForDay(day time.Time) []time.Time {
	var times []time.Time
	for _, t := range ts.times {
		if sameDate(day, t) {
			times = append(times, t)
		}
	}

	sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })

	return times
}

func sameDate(a, b time.Time) bool {
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()

	return aDay == bDay && aMonth == bMonth && aYear == bYear
}
