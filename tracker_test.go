package main

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/roccoblues/tt/test"
)

type testRepo struct {
	data     []byte
	readErr  bool
	writeErr bool
}

func (r *testRepo) Read() ([]byte, error) {
	if r.readErr {
		return nil, errors.New("read failed")
	}
	return r.data, nil
}

func (r *testRepo) Write(d []byte) error {
	if r.writeErr {
		return errors.New("write failed")
	}
	r.data = d
	return nil
}

func Test_newTracker(t *testing.T) {
	tests := []struct {
		name    string
		repo    repository
		want    []*day
		wantErr bool
	}{
		{
			name:    "empty",
			repo:    &testRepo{},
			want:    []*day{},
			wantErr: false,
		},
		{
			name:    "load error",
			repo:    &testRepo{readErr: true},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid json",
			repo:    &testRepo{data: []byte(test.InvalidJSON)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid date",
			repo:    &testRepo{data: []byte(test.InvalidDateJSON)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid time",
			repo:    &testRepo{data: []byte(test.InvalidTimeJSON)},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker, err := newTracker(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("newTracker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				days := tracker.Days()
				if !reflect.DeepEqual(days, tt.want) {
					t.Errorf("newTracker() = %v, want %v", days, tt.want)
				}
			}
		})
	}
}

func Test_tracker_Days(t *testing.T) {
	tests := []struct {
		name string
		json string
		want []*day
	}{
		{
			name: "empty json",
			json: test.EmptyJSON,
			want: []*day{},
		},
		{
			name: "empty day",
			json: test.EmptyDayJSON,
			want: []*day{},
		},
		{
			name: "one day only start",
			json: test.OneDayOnlyStartJSON,
			want: []*day{
				&day{
					Date: test.Time(t, "2018-09-01 10:00"),
					Entries: []*entry{
						&entry{Start: test.Time(t, "2018-09-01 10:00")},
					},
				},
			},
		},
		{
			name: "one day start/end",
			json: test.OneDayStartEndJSON,
			want: []*day{
				&day{
					Date: test.Time(t, "2018-09-01 10:00"),
					Entries: []*entry{
						&entry{Start: test.Time(t, "2018-09-01 10:00"), End: test.Time(t, "2018-09-01 12:00")},
					},
				},
			},
		},
		{
			name: "one day start/end start",
			json: test.OneDayStartEndStartJSON,
			want: []*day{
				&day{
					Date: test.Time(t, "2018-09-01 10:00"),
					Entries: []*entry{
						&entry{Start: test.Time(t, "2018-09-01 10:00"), End: test.Time(t, "2018-09-01 12:00")},
						&entry{Start: test.Time(t, "2018-09-01 13:00")},
					},
				},
			},
		},
		{
			name: "multiple days",
			json: test.MultipleDaysJSON,
			want: []*day{
				&day{
					Date: test.Time(t, "2018-09-01 10:00"),
					Entries: []*entry{
						&entry{Start: test.Time(t, "2018-09-01 10:00"), End: test.Time(t, "2018-09-01 12:00")},
					},
				},
				&day{
					Date: test.Time(t, "2018-09-02 08:00"),
					Entries: []*entry{
						&entry{Start: test.Time(t, "2018-09-02 08:00")},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := newTracker(&testRepo{data: []byte(tt.json)})
			if err != nil {
				t.Fatal(err)
			}
			days := tr.Days()
			if diff := cmp.Diff(tt.want, days); diff != "" {
				t.Errorf("tracker.Days() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_sameDay(t *testing.T) {
	testTime := test.Time(t, "2018-09-01 10:00")

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
	tests := []struct {
		name    string
		before  string
		start   time.Time
		after   string
		wantErr bool
	}{
		{
			name:    "first entry",
			before:  test.EmptyJSON,
			after:   test.OneDayOnlyStartJSON,
			start:   test.Time(t, "2018-09-01 10:00"),
			wantErr: false,
		},
		{
			name:    "second entry",
			before:  test.OneDayStartEndJSON,
			after:   test.OneDayStartEndStartJSON,
			start:   test.Time(t, "2018-09-01 13:00"),
			wantErr: false,
		},
		{
			name:    "new day",
			before:  test.OneDayStartEndJSON,
			after:   test.MultipleDaysJSON,
			start:   test.Time(t, "2018-09-02 08:00"),
			wantErr: false,
		},
		{
			name:    "new day previous not stopped",
			before:  test.OneDayOnlyStartJSON,
			after:   test.MultipleDaysNotEndedJSON,
			start:   test.Time(t, "2018-09-02 08:00"),
			wantErr: false,
		},
		{
			name:    "already started",
			before:  test.OneDayOnlyStartJSON,
			after:   test.OneDayOnlyStartJSON,
			start:   test.Time(t, "2018-09-01 16:00"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &testRepo{data: []byte(tt.before)}
			tr, err := newTracker(repo)
			if err != nil {
				t.Fatal(err)
			}
			err = tr.Start(tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("tracker.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.after, string(repo.data)); diff != "" {
				t.Errorf("tracker.Start() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_tracker_End(t *testing.T) {
	tests := []struct {
		name    string
		before  string
		end     time.Time
		after   string
		wantErr bool
	}{
		{
			name:    "end first entry",
			before:  test.OneDayOnlyStartJSON,
			after:   test.OneDayStartEndJSON,
			end:     test.Time(t, "2018-09-01 12:00"),
			wantErr: false,
		},
		{
			name:    "end second day",
			before:  test.MultipleDaysJSON,
			after:   test.MultipleDaysEndJSON,
			end:     test.Time(t, "2018-09-02 16:00"),
			wantErr: false,
		},
		{
			name:    "not started",
			before:  test.EmptyJSON,
			after:   test.EmptyJSON,
			end:     test.Time(t, "2018-09-01 16:00"),
			wantErr: true,
		},
		{
			name:    "new day previous not stopped",
			before:  test.MultipleDaysNotEndedJSON,
			after:   test.MultipleDaysNotEndedEndJSON,
			end:     test.Time(t, "2018-09-02 16:00"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &testRepo{data: []byte(tt.before)}
			tr, err := newTracker(repo)
			if err != nil {
				t.Fatal(err)
			}
			err = tr.End(tt.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("tracker.End() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.after, string(repo.data)); diff != "" {
				t.Errorf("tracker.End() differs: (-want +got)\n%s", diff)
			}
		})
	}
}
