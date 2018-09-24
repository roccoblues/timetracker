package main

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type timeSheet struct {
	times []time.Time
}

func (t *timeSheet) Start(start time.Time) error {
	c := 0
	for _, tm := range t.times {
		if sameDate(start, tm) {
			c++
		}
	}

	if c%2 != 0 {
		return errors.New("already started")
	}

	t.times = append(t.times, start)

	return nil
}

func (t *timeSheet) End(end time.Time) error {
	c := 0
	for _, tm := range t.times {
		if sameDate(end, tm) {
			c++
		}
	}

	if c%2 == 0 {
		return errors.New("not started")
	}

	t.times = append(t.times, end)

	return nil
}

func (t *timeSheet) Days() []*day {
	days := []*day{}

	for _, t := range t.times {
		if len(days) == 0 {
			days = append(days, newDay(t))
			continue
		}

		day := days[len(days)-1]
		if !sameDate(day.Date, t) {
			days = append(days, newDay(t))
			continue
		}

		last := day.Entries[len(day.Entries)-1]
		if last.End.IsZero() {
			last.End = t
		} else {
			day.Entries = append(day.Entries, &entry{Start: t})
		}
	}

	return days
}

func (t *timeSheet) MarshalJSON() ([]byte, error) {
	data := map[string][]string{}

	for _, tm := range t.times {
		date := tm.Format("2006-01-02")
		if _, exists := data[date]; !exists {
			data[date] = []string{}
		}
		data[date] = append(data[date], tm.Format(timeFormat))
	}

	return json.MarshalIndent(&data, "", "  ")
}

func (t *timeSheet) UnmarshalJSON(bytes []byte) error {
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
			timeString := d + " " + e
			tm, err := time.ParseInLocation(dateTimeFormat, timeString, loc)
			if err != nil {
				return errors.Wrapf(err, "failed to parse entry %s", timeString)
			}
			t.times = append(t.times, tm)
		}
	}

	return nil
}

func sameDate(a, b time.Time) bool {
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()

	return aDay == bDay && aMonth == bMonth && aYear == bYear
}
