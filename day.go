package main

import (
	"time"
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
	return &day{Date: t, Entries: []*entry{&entry{Start: t}}}
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
