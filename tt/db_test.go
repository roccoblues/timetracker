package main

import (
	"testing"
	"time"
)

func Test_sameDay(t *testing.T) {
	testTime, err := time.Parse("2006-01-02 15:04", "2018-09-11 13:40")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		a    time.Time
		b    time.Time
		want bool
	}{
		{
			name: "same day",
			a:    testTime,
			b:    testTime.Add(time.Minute * 5),
			want: true,
		},
		{
			name: "next day",
			a:    testTime,
			b:    testTime.Add(time.Hour * 25),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sameDay(tt.a, tt.b); got != tt.want {
				t.Errorf("SameDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_worked(t *testing.T) {
	tests := []struct {
		name     string
		times    []string
		duration time.Duration
	}{
		{
			name:     "empty",
			times:    []string{},
			duration: 0,
		},
		{
			name:     "only start",
			times:    []string{"08:00"},
			duration: 0,
		},
		{
			name:     "one entry",
			times:    []string{"08:00", "12:00"},
			duration: time.Duration(4) * time.Hour,
		},
		{
			name:     "one entry with trailing start",
			times:    []string{"08:00", "12:00", "13:00"},
			duration: time.Duration(4) * time.Hour,
		},
		{
			name:     "multiple entries",
			times:    []string{"08:00", "12:00", "13:00", "15:30"},
			duration: time.Duration(390) * time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			times := []time.Time{}
			for _, s := range tt.times {
				e, err := time.Parse("15:04", s)
				if err != nil {
					t.Fatal(err)
				}
				times = append(times, e)
			}

			if got := worked(times); got != tt.duration {
				t.Errorf("worked() = %v, want %v", got, tt.duration)
			}
		})
	}
}
