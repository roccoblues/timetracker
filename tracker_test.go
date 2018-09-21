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
	times    []time.Time
	readErr  bool
	writeErr bool
}

func (r *testRepo) Read() ([]time.Time, error) {
	if r.readErr {
		return nil, errors.New("read failed")
	}
	return r.times, nil
}

func (r *testRepo) Write(times []time.Time) error {
	if r.writeErr {
		return errors.New("write failed")
	}
	r.times = times
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
			name:    "load error",
			repo:    &testRepo{readErr: true},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			repo:    &testRepo{},
			want:    []*day{},
			wantErr: false,
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

func Test_tracker_Start(t *testing.T) {
	tests := []struct {
		name         string
		start        time.Time
		before       []time.Time
		after        []time.Time
		wantErr      bool
		repoWriteErr bool
	}{
		{
			name:   "first entry",
			start:  test.Time(t, "2018-09-01 10:00"),
			before: []time.Time{},
			after: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
			},
			wantErr: false,
		},
		{
			name:  "second entry",
			start: test.Time(t, "2018-09-01 13:00"),
			before: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-01 12:00"),
			},
			after: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-01 13:00"),
			},
			wantErr: false,
		},
		{
			name:  "new day",
			start: test.Time(t, "2018-09-02 08:00"),
			before: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-01 12:00"),
			},
			after: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-02 08:00"),
			},
			wantErr: false,
		},
		{
			name:  "new day previous not stopped",
			start: test.Time(t, "2018-09-02 08:00"),
			before: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
			},
			after: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-02 08:00"),
			},
			wantErr: false,
		},
		{
			name:  "already started",
			start: test.Time(t, "2018-09-01 16:00"),
			before: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
			},
			after: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
			},
			wantErr: true,
		},
		{
			name:         "write failure",
			start:        test.Time(t, "2018-09-01 10:00"),
			before:       []time.Time{},
			after:        []time.Time{},
			wantErr:      true,
			repoWriteErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &testRepo{times: tt.before, writeErr: tt.repoWriteErr}
			tr, err := newTracker(repo)
			if err != nil {
				t.Fatal(err)
			}
			err = tr.Start(tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("tracker.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := repo.Read()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.after, got); diff != "" {
				t.Errorf("tracker.Start() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_tracker_End(t *testing.T) {
	tests := []struct {
		name         string
		end          time.Time
		before       []time.Time
		after        []time.Time
		wantErr      bool
		repoWriteErr bool
	}{
		{
			name: "end first entry",
			end:  test.Time(t, "2018-09-01 12:00"),
			before: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
			},
			after: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-01 12:00"),
			},
			wantErr: false,
		},
		{
			name: "end second day",
			end:  test.Time(t, "2018-09-02 16:00"),
			before: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-02 10:00"),
			},
			after: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-02 10:00"),
				test.Time(t, "2018-09-02 16:00"),
			},
			wantErr: false,
		},
		{
			name:    "not started",
			end:     test.Time(t, "2018-09-01 16:00"),
			before:  []time.Time{},
			after:   []time.Time{},
			wantErr: true,
		},
		{
			name: "new day previous not stopped",
			end:  test.Time(t, "2018-09-02 16:00"),
			before: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-02 10:00"),
			},
			after: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
				test.Time(t, "2018-09-02 10:00"),
				test.Time(t, "2018-09-02 16:00"),
			},
			wantErr: false,
		},
		{
			name: "write failure",
			end:  test.Time(t, "2018-09-01 12:00"),
			before: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
			},
			after: []time.Time{
				test.Time(t, "2018-09-01 08:00"),
			},
			wantErr:      true,
			repoWriteErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &testRepo{times: tt.before, writeErr: tt.repoWriteErr}
			tr, err := newTracker(repo)
			if err != nil {
				t.Fatal(err)
			}
			err = tr.End(tt.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("tracker.End() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := repo.Read()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.after, got); diff != "" {
				t.Errorf("tracker.End() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_tracker_Days(t *testing.T) {
	tests := []struct {
		name  string
		times []time.Time
		want  []*day
	}{
		{
			name:  "empty",
			times: []time.Time{},
			want:  []*day{},
		},
		{
			name: "one day only start",
			times: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
			},
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
			times: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
				test.Time(t, "2018-09-01 12:00"),
			},
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
			times: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-01 13:00"),
			},
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
			times: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-02 08:00"),
			},
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
			tr, err := newTracker(&testRepo{times: tt.times})
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

func Test_sameDate(t *testing.T) {
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
			if got := sameDate(tt.a, tt.b); got != tt.want {
				t.Errorf("sameDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
