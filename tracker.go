package main

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type repository interface {
	Read() ([]byte, error)
	Write([]byte) error
}

type tracker struct {
	repo  repository
	loc   *time.Location
	times []time.Time
}

func newTracker(repo repository) (*tracker, error) {
	tracker := tracker{
		repo: repo,
		loc:  time.Now().Location(),
	}

	if err := tracker.load(); err != nil {
		return nil, errors.Wrap(err, "failed to load data")
	}

	return &tracker, nil
}

func (t *tracker) Start(start time.Time) error {
	count := 0
	for _, tm := range t.times {
		if sameDay(start, tm) {
			count++
		}
	}

	if count%2 != 0 {
		return errors.New("already started")
	}

	t.times = append(t.times, start)

	if err := t.save(); err != nil {
		return errors.Wrap(err, "save failed")
	}

	return nil
}

func (t *tracker) End(end time.Time) error {
	count := 0
	for _, tm := range t.times {
		if sameDay(end, tm) {
			count++
		}
	}

	if count%2 == 0 {
		return errors.New("not started")
	}

	t.times = append(t.times, end)

	if err := t.save(); err != nil {
		return errors.Wrap(err, "save failed")
	}

	return nil
}

func (t *tracker) Days() []*day {
	days := []*day{}
	for _, t := range t.times {
		if len(days) == 0 {
			days = append(days, newDay(t))
			continue
		}

		day := days[len(days)-1]
		if !sameDay(day.Date, t) {
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

func (t *tracker) load() error {
	data, err := t.repo.Read()
	if err != nil {
		return errors.Wrap(err, "failed to read data")
	}

	if len(data) > 0 {
		if err := t.decode(data); err != nil {
			return errors.Wrap(err, "decode failed")
		}
	}

	return nil
}

func (t *tracker) save() error {
	data, err := t.encode()
	if err != nil {
		return errors.Wrap(err, "encode failed")
	}

	if err := t.repo.Write(data); err != nil {
		return errors.Wrap(err, "failed to write data")
	}

	return nil
}

func (t *tracker) decode(data []byte) error {
	var decoded map[string][]string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return errors.Wrap(err, "json decode failed")
	}

	dates := make([]string, 0, len(decoded))
	for d := range decoded {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	for _, d := range dates {
		for _, e := range decoded[d] {
			timeString := d + " " + e
			tm, err := time.ParseInLocation("2006-01-02 15:04", timeString, t.loc)
			if err != nil {
				return errors.Wrapf(err, "failed to parse entry %s", timeString)
			}
			t.times = append(t.times, tm)
		}
	}

	return nil
}

func (t *tracker) encode() ([]byte, error) {
	data := map[string][]string{}

	for _, t := range t.times {
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

func sameDay(a, b time.Time) bool {
	return a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year()
}
