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

type jsonFile struct {
	path string
}

func newJSONFile(path string) *jsonFile {
	return &jsonFile{path: path}
}

func (j *jsonFile) Read() ([]*day, error) {
	json := []byte(emptyJSON)

	if _, err := os.Stat(j.path); err == nil {
		json, err = ioutil.ReadFile(j.path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read file '%s'", j.path)
		}
	}

	days, err := decode(json)
	if err != nil {
		return nil, errors.Wrap(err, "decode failed")
	}

	return days, nil
}

func (j *jsonFile) Write(days []*day) error {
	json, err := encode(days)
	if err != nil {
		return errors.Wrap(err, "encode failed")
	}

	if err := ioutil.WriteFile(j.path, json, 0644); err != nil {
		return errors.Wrapf(err, "failed to write to '%s'", j.path, err)
	}

	return nil
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

		for i, e := range data[d] {
			t, err := time.Parse("2006-01-02 15:04", d+" "+e)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse time %s", e)
			}
			if i%2 == 0 {
				day.StartEntry(t)
			} else {
				day.EndEntry(t)
			}
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
