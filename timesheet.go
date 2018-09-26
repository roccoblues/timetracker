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
	c := 0
	for _, t := range ts.times {
		if sameDate(start, t) {
			c++
		}
	}

	if c%2 != 0 {
		return fmt.Errorf("already started")
	}

	ts.times = append(ts.times, start)

	return nil
}

func (ts *timeSheet) End(end time.Time) error {
	c := 0
	for _, t := range ts.times {
		if sameDate(end, t) {
			c++
		}
	}

	if c%2 == 0 {
		return fmt.Errorf("not started")
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

func (ts *timeSheet) Print(out io.Writer, roundTo time.Duration) {
	var times []time.Time
	var week int
	var day time.Time

	for n, t := range ts.times {
		if day.IsZero() || sameDate(day, t) {
			times = append(times, t)
			day = t
			if n != len(ts.times)-1 {
				continue
			}
		}

		var hours time.Duration
		var start time.Time
		for i, t := range times {
			if i%2 == 0 {
				start = t
			} else {
				hours += t.Sub(start)
			}
		}

		fmt.Fprintf(out, "%s  %.2f ", day.Format(dateFormat), hours.Round(roundTo).Hours())

		for i, t := range times {
			if i%2 == 0 {
				fmt.Fprintf(out, " %s-", t.Format(timeFormat))
			} else {
				fmt.Fprintf(out, "%s", t.Format(timeFormat))
			}
		}

		fmt.Fprintln(out, "")

		_, w := t.ISOWeek()
		if week > 0 && week != w {
			fmt.Fprintln(out, "")
		}

		day = t
		week = w
		times = []time.Time{t}
		hours = 0
	}
}

func sameDate(a, b time.Time) bool {
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()

	return aDay == bDay && aMonth == bMonth && aYear == bYear
}
