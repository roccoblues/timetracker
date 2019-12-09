package timesheet

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"
)

type dateTimes map[string][]string

func unmarshal(r io.Reader, dateFormat, timeFormat string) ([]time.Time, error) {
	var dt dateTimes

	dec := json.NewDecoder(r)
	err := dec.Decode(&dt)
	if err != nil && err != io.EOF {
		return nil, err
	}

	dateTimeFormat := fmt.Sprintf("%s %s", dateFormat, timeFormat)

	loc := time.Now().Location()

	var times []time.Time
	for dateStr, timeStrs := range dt {
		for _, t := range timeStrs {
			dateTime := fmt.Sprintf("%s %s", dateStr, t)
			tm, err := time.ParseInLocation(dateTimeFormat, dateTime, loc)
			if err != nil {
				return nil, err
			}
			times = append(times, tm)
		}
	}

	sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })

	return times, nil
}

func marshal(w io.Writer, times []time.Time, dateFormat, timeFormat string) error {
	dt := dateTimes{}

	for _, t := range times {
		date := t.Format(dateFormat)
		if _, exists := dt[date]; !exists {
			dt[date] = []string{}
		}
		dt[date] = append(dt[date], t.Format(timeFormat))
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc.Encode(dt)
}
