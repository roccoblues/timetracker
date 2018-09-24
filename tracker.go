package main

import (
	"time"

	"github.com/pkg/errors"
)

type repository interface {
	Read() ([]time.Time, error)
	Write([]time.Time) error
}

type tracker struct {
	repo  repository
	times []time.Time
}

func newTracker(repo repository) (*tracker, error) {
	times, err := repo.Read()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read data")
	}

	tracker := tracker{
		repo:  repo,
		times: times,
	}

	return &tracker, nil
}

func (t *tracker) Start(start time.Time) error {
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

	if err := t.repo.Write(t.times); err != nil {
		return errors.Wrap(err, "failed to write data")
	}

	return nil
}

func (t *tracker) End(end time.Time) error {
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

	if err := t.repo.Write(t.times); err != nil {
		return errors.Wrap(err, "failed to write data")
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

func sameDate(a, b time.Time) bool {
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()

	return aDay == bDay && aMonth == bMonth && aYear == bYear
}
