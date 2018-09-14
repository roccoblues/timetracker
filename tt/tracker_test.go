package main

import (
	"reflect"
	"testing"
	"time"
)

type testStorage struct {
	times []time.Time
}

func (s *testStorage) Read() ([]time.Time, error) {
	return s.times, nil
}
func (s *testStorage) Write(times []time.Time) error {
	s.times = times
	return nil
}

func Test_newTracker(t *testing.T) {
	storage := &testStorage{}

	tests := []struct {
		name  string
		input persistence
		want  *tracker
	}{
		{
			name:  "returns a valid tracker",
			input: storage,
			want:  &tracker{db: storage},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTracker(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTracker() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
