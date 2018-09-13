package main

import (
	"time"

	"github.com/pkg/errors"
)

type day struct {
	Date    time.Time
	Entries []*entry
}

type entry struct {
	Start time.Time
	End   time.Time
}

func newDay(t time.Time) *day {
	return &day{Date: t}
}

func (d *day) StartEntry(t time.Time) error {
	if len(d.Entries) == 0 {
		d.Entries = append(d.Entries, &entry{Start: t})
		return nil
	}

	e := d.Entries[len(d.Entries)-1]

	if !e.Start.IsZero() && e.End.IsZero() {
		return errors.New("already started")
	}

	d.Entries = append(d.Entries, &entry{Start: t})
	return nil
}

func (d *day) StopEntry(t time.Time) error {
	if len(d.Entries) == 0 {
		return errors.New("not started")
	}

	e := d.Entries[len(d.Entries)-1]

	if !e.End.IsZero() {
		return errors.New("not started")
	}

	e.End = t
	return nil
}

func (d *day) AddTime(t time.Time) {
	if len(d.Entries) == 0 {
		d.Entries = append(d.Entries, &entry{Start: t})
		return
	}

	e := d.Entries[len(d.Entries)-1]

	if e.Start.IsZero() {
		e.Start = t
	} else if e.End.IsZero() {
		e.End = t
	} else {
		d.Entries = append(d.Entries, &entry{Start: t})
	}
}

func (d *day) Time() time.Duration {
	var worked time.Duration

	for _, e := range d.Entries {
		if !e.Start.IsZero() && !e.End.IsZero() {
			worked += e.End.Sub(e.Start)
		}
	}

	return worked
}
