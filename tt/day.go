package main

import (
	"time"

	"github.com/pkg/errors"
)

type day struct {
	date  time.Time
	times []time.Time
}

func (d *day) addStartTime(t time.Time) error {
	if len(d.times)%2 != 0 {
		return errors.New("already started")
	}

	d.times = append(d.times, t)
	return nil
}

func (d *day) addStopTime(t time.Time) error {
	if len(d.times)%2 == 0 {
		return errors.New("not started")
	}

	d.times = append(d.times, t)
	return nil
}
