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

func (j *jsonFile) Read() ([]time.Time, error) {
	json := []byte(emptyJSON)

	if _, err := os.Stat(j.path); err == nil {
		json, err = ioutil.ReadFile(j.path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read file '%s'", j.path)
		}
	}

	times, err := decode(json)
	if err != nil {
		return nil, errors.Wrap(err, "decode failed")
	}

	return times, nil
}

func (j *jsonFile) Write(times []time.Time) error {
	json, err := encode(times)
	if err != nil {
		return errors.Wrap(err, "encode failed")
	}

	if err := ioutil.WriteFile(j.path, json, 0644); err != nil {
		return errors.Wrapf(err, "failed to write to '%s'", j.path, err)
	}

	return nil
}

func decode(input []byte) ([]time.Time, error) {
	var data map[string][]string
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, errors.Wrap(err, "json decode failed")
	}

	dates := make([]string, 0, len(data))
	for d := range data {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	times := []time.Time{}
	for _, d := range dates {
		for _, e := range data[d] {
			timeString := d + " " + e
			t, err := time.Parse("2006-01-02 15:04", timeString)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse entry %s", timeString)
			}
			times = append(times, t)
		}
	}

	return times, nil
}

func encode(times []time.Time) ([]byte, error) {
	data := map[string][]string{}

	for _, t := range times {
		date := t.Format("2006-01-02")
		if _, exists := data[date]; !exists {
			data[date] = []string{}
		}
		data[date] = append(data[date], t.Format("15:04"))
	}

	encoded, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "json encode failed")
	}

	return encoded, nil
}
