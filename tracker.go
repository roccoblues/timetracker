package main

import (
	"time"

	"github.com/pkg/errors"
)

type persistence interface {
	Read() ([]time.Time, error)
	Write([]time.Time) error
}

type tracker struct {
	db persistence
}

func newTracker(db persistence) *tracker {
	return &tracker{db: db}
}

func (t *tracker) Start(start time.Time) error {
	times, err := t.db.Read()
	if err != nil {
		return errors.Wrap(err, "read failed")
	}

	if len(times)%2 != 0 {
		return errors.New("already started")
	}

	times = append(times, start)

	err = t.db.Write(times)
	if err != nil {
		return errors.Wrap(err, "write failed")
	}

	return nil
}

func (t *tracker) End(end time.Time) error {
	times, err := t.db.Read()
	if err != nil {
		return errors.Wrap(err, "read failed")
	}

	if len(times)%2 == 0 {
		return errors.New("not started")
	}

	times = append(times, end)

	err = t.db.Write(times)
	if err != nil {
		return errors.Wrap(err, "write failed")
	}

	return nil
}

func (t *tracker) Days() ([]*day, error) {
	times, err := t.db.Read()
	if err != nil {
		return nil, errors.Wrap(err, "read failed")
	}

	days := []*day{}
	for _, t := range times {
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
	return days, nil
}

func sameDay(a, b time.Time) bool {
	return a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year()
}
