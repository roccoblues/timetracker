package main

import (
	"time"

	"github.com/pkg/errors"
)

type persistence interface {
	Read() ([]*day, error)
	Write([]*day) error
}

type tracker struct {
	db persistence
}

func newTracker(db persistence) *tracker {
	return &tracker{db: db}
}

func (t *tracker) Start(start time.Time) error {
	days, err := t.db.Read()
	if err != nil {
		return errors.Wrap(err, "read failed")
	}

	day := days[len(days)-1]

	if day == nil || !sameDay(day.Date, start) {
		day = newDay(start)
		days = append(days, day)
	}

	err = day.StartEntry(start)
	if err != nil {
		return errors.Wrap(err, "adding start time to day failed")
	}

	err = t.db.Write(days)
	if err != nil {
		return errors.Wrap(err, "write failed")
	}

	return nil
}

func (t *tracker) End(end time.Time) error {
	days, err := t.db.Read()
	if err != nil {
		return errors.Wrap(err, "read failed")
	}

	day := days[len(days)-1]

	if err := day.EndEntry(end); err != nil {
		return errors.Wrap(err, "adding stop time to day failed")
	}

	err = t.db.Write(days)
	if err != nil {
		return errors.Wrap(err, "write failed")
	}

	return nil
}

func (t *tracker) All() ([]*day, error) {
	return t.db.Read()
}

func sameDay(a, b time.Time) bool {
	return a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year()
}
