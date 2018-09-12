package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type db struct {
	days     []*day
	Modified bool
}

func (db *db) decode(input []byte) error {
	var data map[string][]string
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "json decode failed")
	}

	dates := make([]string, 0, len(data))
	for d := range data {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	days := []*day{}
	for _, d := range dates {
		date, err := time.Parse("2006-01-02", d)
		if err != nil {
			return errors.Wrapf(err, "failed to parse date %s", d)
		}

		times := []time.Time{}
		for _, e := range data[d] {
			t, err := time.Parse("2006-01-02 15:04", d+" "+e)
			if err != nil {
				return errors.Wrapf(err, "failed to parse time %s", e)
			}
			times = append(times, t)
		}

		days = append(days, &day{date: date, times: times})
	}

	db.days = days

	return nil
}

func (db *db) encode() ([]byte, error) {
	data := map[string][]string{}

	for _, day := range db.days {
		times := []string{}
		for _, t := range day.times {
			times = append(times, t.Format("15:04"))
		}
		data[day.date.Format("2006-01-02")] = times
	}

	encoded, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "json encode failed")
	}

	return encoded, nil
}

func (db *db) print(out io.Writer) {
	var week int
	for _, day := range db.days {
		_, w := day.date.ISOWeek()
		if week > 0 && week != w {
			fmt.Fprintln(out, "")
		}
		week = w

		fmt.Fprintf(out, "%s  %.2f  ", day.date.Format("02.01.2006"), worked(day.times).Round(roundTo).Hours())

		for i, t := range day.times {
			fmt.Fprintf(out, "%s", t.Round(roundTo).Format("15:04"))
			if i%2 == 0 {
				fmt.Fprint(out, "-")
			} else {
				fmt.Fprint(out, " ")
			}
			if i == len(day.times)-1 {
				fmt.Fprintln(out, "")
			}
		}
	}
}

func (db *db) addStartTime(t time.Time) error {
	for _, d := range db.days {
		if sameDay(d.date, t) {
			if err := d.addStartTime(t); err != nil {
				return err
			}
			db.Modified = true
		}
	}

	db.days = append(db.days, &day{date: t, times: []time.Time{t}})
	db.Modified = true

	return nil
}

func (db *db) addStopTime(t time.Time) error {
	for _, d := range db.days {
		if sameDay(d.date, t) {
			return d.addStopTime(t)
		}
	}

	return errors.New("not started")
}

func sameDay(a, b time.Time) bool {
	return a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year()
}

func worked(times []time.Time) time.Duration {
	var worked time.Duration
	var start time.Time

	for i, t := range times {
		if i%2 == 0 {
			start = t
		} else {
			worked += t.Sub(start)
		}
	}

	return worked
}
