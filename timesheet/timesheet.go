package timesheet

import (
	"fmt"
	"io"
	"time"
)

type Sheet struct {
	DateFormat string
	TimeFormat string
	Times      []time.Time
}

func Load(r io.Reader, dateFormat, timeFormat string) (*Sheet, error) {
	times, err := unmarshal(r, dateFormat, timeFormat)
	if err != nil {
		return nil, err
	}

	sheet := &Sheet{
		DateFormat: dateFormat,
		TimeFormat: timeFormat,
		Times:      times,
	}

	return sheet, nil
}

func (s *Sheet) Save(w io.Writer) error {
	return marshal(w, s.Times, s.DateFormat, s.TimeFormat)
}

func (s *Sheet) Start(start time.Time) error {
	var last time.Time
	c := 0
	for _, t := range s.Times {
		if sameDate(start, t) {
			c++
			last = t
		}
	}

	if c%2 != 0 {
		return fmt.Errorf("already started")
	}

	if start.Before(last) {
		return fmt.Errorf("start time %s is ealier as last end time %s", start.Format(s.TimeFormat), last.Format(s.TimeFormat))
	}

	s.Times = append(s.Times, start)

	return nil
}

func (s *Sheet) End(end time.Time) error {
	var last time.Time
	c := 0
	for _, t := range s.Times {
		if sameDate(end, t) {
			c++
			last = t
		}
	}

	if c%2 == 0 {
		return fmt.Errorf("not started")
	}

	if end.Before(last) {
		return fmt.Errorf("end time %s is earlier as last start time %s", end.Format(s.TimeFormat), last.Format(s.TimeFormat))
	}

	s.Times = append(s.Times, end)

	return nil
}

func (s *Sheet) Print(roundTo time.Duration, w io.Writer) {
	print(s.Times, roundTo, s.DateFormat, s.TimeFormat, w)
}

func (s *Sheet) PrintMonth(month time.Month, roundTo time.Duration, w io.Writer) {
	var times []time.Time

	for _, t := range s.Times {
		if t.Month() != month {
			continue
		}
		times = append(times, t)
	}

	print(times, roundTo, s.DateFormat, s.TimeFormat, w)
}

func sameDate(a, b time.Time) bool {
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()

	return aDay == bDay && aMonth == bMonth && aYear == bYear
}
