package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/pkg/errors"
)

const emptyJSON = "{}"

type db struct {
	Path string
}

func newDB(path string) *db {
	return &db{Path: path}
}

func (db *db) AddStart(t time.Time) error {
	days, err := read(db.Path)
	if err != nil {
		return errors.Wrap(err, "read failed")
	}

	days, err = addStart(days, t)
	if err != nil {
		return errors.Wrap(err, "failed to add start time")
	}

	err = write(days, db.Path)
	if err != nil {
		return errors.Wrap(err, "write failed")
	}

	return nil
}

func (db *db) AddStop(t time.Time) error {
	days, err := read(db.Path)
	if err != nil {
		return errors.Wrap(err, "read failed")
	}

	days, err = addStop(days, t)
	if err != nil {
		return errors.Wrap(err, "failed to add stop time")
	}

	err = write(days, db.Path)
	if err != nil {
		return errors.Wrap(err, "write failed")
	}

	return nil
}

func (db *db) All() ([]*day, error) {
	days, err := read(db.Path)
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	return days, nil
}

func decode(input []byte) ([]*day, error) {
	var data map[string][]string
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, errors.Wrap(err, "json decode failed")
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
			return nil, errors.Wrapf(err, "failed to parse date %s", d)
		}

		day := newDay(date)

		for _, e := range data[d] {
			t, err := time.Parse("2006-01-02 15:04", d+" "+e)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse time %s", e)
			}
			day.AddTime(t)
		}

		days = append(days, day)
	}

	return days, nil
}

func encode(days []*day) ([]byte, error) {
	data := map[string][]string{}

	for _, day := range days {
		times := []string{}
		for _, e := range day.Entries {
			times = append(times, e.Start.Format("15:04"), e.End.Format("15:04"))
		}
		data[day.Date.Format("2006-01-02")] = times
	}

	encoded, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "json encode failed")
	}

	return encoded, nil
}

func read(path string) ([]*day, error) {
	json := []byte(emptyJSON)

	if _, err := os.Stat(path); err == nil {
		json, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read file '%s'", path)
		}
	}

	days, err := decode(json)
	if err != nil {
		return nil, errors.Wrap(err, "decode failed")
	}

	return days, nil
}

func write(days []*day, path string) error {
	json, err := encode(days)
	if err != nil {
		return errors.Wrap(err, "encode failed")
	}

	if err := ioutil.WriteFile(path, json, 0644); err != nil {
		return errors.Wrapf(err, "failed to write to '%s'", path, err)
	}

	return nil
}

func addStart(days []*day, t time.Time) ([]*day, error) {
	day := days[len(days)-1]

	if day == nil || !sameDay(day.Date, t) {
		day = newDay(t)
		days = append(days, day)
	}

	err := day.StartEntry(t)
	if err != nil {
		return nil, errors.Wrap(err, "adding start time to day failed")
	}

	return days, nil
}

func addStop(days []*day, t time.Time) ([]*day, error) {
	day := days[len(days)-1]

	if err := day.StopEntry(t); err != nil {
		return nil, errors.Wrap(err, "adding stop time to day failed")
	}

	return days, nil
}

func sameDay(a, b time.Time) bool {
	return a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year()
}
