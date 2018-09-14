package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_newDay(t *testing.T) {
	testTime := time.Now()

	tests := []struct {
		name string
		t    time.Time
		want *day
	}{
		{
			name: "returns valid day",
			t:    testTime,
			want: &day{Date: testTime, Entries: []*entry{&entry{Start: testTime}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newDay(tt.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_day_Time(t *testing.T) {
	testTime := time.Now()

	tests := []struct {
		name    string
		Entries []*entry
		want    time.Duration
	}{
		{
			name:    "empty",
			Entries: []*entry{},
			want:    0,
		},
		{
			name: "only start",
			Entries: []*entry{
				&entry{
					Start: testTime.Add(time.Hour * -5),
				},
			},
			want: 0,
		},
		{
			name: "one entry",
			Entries: []*entry{
				&entry{
					Start: testTime.Add(time.Hour * -5),
					End:   testTime,
				},
			},
			want: time.Duration(5) * time.Hour,
		},
		{
			name: "one entry and trailing start",
			Entries: []*entry{
				&entry{
					Start: testTime.Add(time.Hour * -5),
					End:   testTime.Add(time.Hour * -2),
				},
				&entry{
					Start: testTime.Add(time.Hour * -1),
				},
			},
			want: time.Duration(3) * time.Hour,
		},
		{
			name: "two entries",
			Entries: []*entry{
				&entry{
					Start: testTime.Add(time.Hour * -5),
					End:   testTime.Add(time.Hour * -2),
				},
				&entry{
					Start: testTime.Add(time.Hour * -1),
					End:   testTime,
				},
			},
			want: time.Duration(4) * time.Hour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &day{Entries: tt.Entries}
			if got := d.Time(); got != tt.want {
				t.Errorf("day.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}
