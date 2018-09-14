package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
)

type testStorage struct {
	times      []time.Time
	readError  bool
	writeError bool
}

func (s *testStorage) Read() ([]time.Time, error) {
	if s.readError {
		return nil, errors.New("")
	}
	return s.times, nil
}
func (s *testStorage) Write(times []time.Time) error {
	if s.writeError {
		return errors.New("")
	}
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

func Test_tracker_Start(t *testing.T) {
	start := time.Now()

	tests := []struct {
		name           string
		start          time.Time
		initialStorage persistence
		wantStorage    persistence
		wantErr        bool
	}{
		{
			name:           "works",
			start:          start,
			initialStorage: &testStorage{},
			wantStorage:    &testStorage{times: []time.Time{start}},
			wantErr:        false,
		},
		{
			name:           "already started",
			start:          start,
			initialStorage: &testStorage{times: []time.Time{start.Add(time.Hour * -1)}},
			wantStorage:    &testStorage{times: []time.Time{start.Add(time.Hour * -1)}},
			wantErr:        true,
		},
		{
			name:           "read failed",
			start:          start,
			initialStorage: &testStorage{readError: true},
			wantStorage:    &testStorage{readError: true},
			wantErr:        true,
		},
		{
			name:           "write failed",
			start:          start,
			initialStorage: &testStorage{writeError: true},
			wantStorage:    &testStorage{writeError: true},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := newTracker(tt.initialStorage)
			if err := tracker.Start(tt.start); (err != nil) != tt.wantErr {
				t.Errorf("tracker.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tracker.db, tt.wantStorage) {
				t.Errorf("encode() = %s, want %s", tracker.db, tt.wantStorage)
			}
		})
	}
}

func Test_tracker_End(t *testing.T) {
	start := time.Now().Add(time.Hour * -7)
	end := time.Now()

	tests := []struct {
		name           string
		end            time.Time
		initialStorage persistence
		wantStorage    persistence
		wantErr        bool
	}{
		{
			name:           "works",
			end:            end,
			initialStorage: &testStorage{times: []time.Time{start}},
			wantStorage:    &testStorage{times: []time.Time{start, end}},
			wantErr:        false,
		},
		{
			name:           "not",
			end:            end,
			initialStorage: &testStorage{},
			wantStorage:    &testStorage{},
			wantErr:        true,
		},
		{
			name:           "read failed",
			end:            end,
			initialStorage: &testStorage{readError: true},
			wantStorage:    &testStorage{readError: true},
			wantErr:        true,
		},
		{
			name:           "write failed",
			end:            end,
			initialStorage: &testStorage{times: []time.Time{start}, writeError: true},
			wantStorage:    &testStorage{times: []time.Time{start}, writeError: true},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := newTracker(tt.initialStorage)
			if err := tracker.End(tt.end); (err != nil) != tt.wantErr {
				t.Errorf("tracker.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tracker.db, tt.wantStorage) {
				t.Errorf("encode() = %s, want %s", tracker.db, tt.wantStorage)
			}
		})
	}
}

func Test_tracker_Days(t *testing.T) {
	testTime1, err := time.Parse("2006-01-02 15:04", "2018-09-01 10:00")
	if err != nil {
		t.Fatal(err)
	}
	testTime2 := testTime1.Add(time.Hour * 2)
	testTime3 := testTime1.Add(time.Hour * 4)
	testTime4 := testTime1.Add(time.Hour * 24)

	tests := []struct {
		name    string
		storage persistence
		want    []*day
		wantErr bool
	}{
		{
			name:    "read failed",
			storage: &testStorage{readError: true},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "works",
			storage: &testStorage{times: []time.Time{testTime1, testTime2, testTime3, testTime4}},
			want: []*day{
				&day{
					Date: testTime1,
					Entries: []*entry{
						&entry{
							Start: testTime1,
							End:   testTime2,
						},
						&entry{
							Start: testTime3,
						},
					},
				},
				&day{
					Date: testTime4,
					Entries: []*entry{
						&entry{
							Start: testTime4,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := newTracker(tt.storage)
			got, err := tracker.Days()
			if (err != nil) != tt.wantErr {
				t.Errorf("tracker.Days() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tracker.Days() = %v, want %v", got, tt.want)
			}
		})
	}
}
