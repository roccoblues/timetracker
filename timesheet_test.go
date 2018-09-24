package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func readFile(t *testing.T, path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func newTime(t *testing.T, str string) time.Time {
	tm, err := time.ParseInLocation(dateTimeFormat, str, time.Now().Location())
	if err != nil {
		t.Fatal(err)
	}
	return tm
}

func Test_timeSheet_sameDate(t *testing.T) {
	testTime := newTime(t, "2018-09-01 10:00")

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

func Test_timeSheet_Days(t *testing.T) {
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
				newTime(t, "2018-09-01 10:00"),
			},
			want: []*day{
				&day{
					Date: newTime(t, "2018-09-01 10:00"),
					Entries: []*entry{
						&entry{Start: newTime(t, "2018-09-01 10:00")},
					},
				},
			},
		},
		{
			name: "one day start/end",
			times: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
			},
			want: []*day{
				&day{
					Date: newTime(t, "2018-09-01 10:00"),
					Entries: []*entry{
						&entry{Start: newTime(t, "2018-09-01 10:00"), End: newTime(t, "2018-09-01 12:00")},
					},
				},
			},
		},
		{
			name: "one day start/end start",
			times: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-01 13:00"),
			},
			want: []*day{
				&day{
					Date: newTime(t, "2018-09-01 10:00"),
					Entries: []*entry{
						&entry{Start: newTime(t, "2018-09-01 10:00"), End: newTime(t, "2018-09-01 12:00")},
						&entry{Start: newTime(t, "2018-09-01 13:00")},
					},
				},
			},
		},
		{
			name: "multiple days",
			times: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-02 08:00"),
			},
			want: []*day{
				&day{
					Date: newTime(t, "2018-09-01 10:00"),
					Entries: []*entry{
						&entry{Start: newTime(t, "2018-09-01 10:00"), End: newTime(t, "2018-09-01 12:00")},
					},
				},
				&day{
					Date: newTime(t, "2018-09-02 08:00"),
					Entries: []*entry{
						&entry{Start: newTime(t, "2018-09-02 08:00")},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &timeSheet{times: tt.times}
			days := ts.Days()
			if diff := cmp.Diff(tt.want, days); diff != "" {
				t.Errorf("timeSheet.Days() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_timeSheet_Start(t *testing.T) {
	tests := []struct {
		name    string
		start   time.Time
		before  []time.Time
		after   []time.Time
		wantErr bool
	}{
		{
			name:   "first entry",
			start:  newTime(t, "2018-09-01 10:00"),
			before: []time.Time{},
			after: []time.Time{
				newTime(t, "2018-09-01 10:00"),
			},
			wantErr: false,
		},
		{
			name:  "second entry",
			start: newTime(t, "2018-09-01 13:00"),
			before: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-01 12:00"),
			},
			after: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-01 13:00"),
			},
			wantErr: false,
		},
		{
			name:  "new day",
			start: newTime(t, "2018-09-02 08:00"),
			before: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-01 12:00"),
			},
			after: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-02 08:00"),
			},
			wantErr: false,
		},
		{
			name:  "new day previous not stopped",
			start: newTime(t, "2018-09-02 08:00"),
			before: []time.Time{
				newTime(t, "2018-09-01 08:00"),
			},
			after: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-02 08:00"),
			},
			wantErr: false,
		},
		{
			name:  "already started",
			start: newTime(t, "2018-09-01 16:00"),
			before: []time.Time{
				newTime(t, "2018-09-01 08:00"),
			},
			after: []time.Time{
				newTime(t, "2018-09-01 08:00"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &timeSheet{times: tt.before}
			err := ts.Start(tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeSheet.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.after, ts.times); diff != "" {
				t.Errorf("timeSheet.Start() times differ: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_timeSheet_End(t *testing.T) {
	tests := []struct {
		name    string
		end     time.Time
		before  []time.Time
		after   []time.Time
		wantErr bool
	}{
		{
			name: "end first entry",
			end:  newTime(t, "2018-09-01 12:00"),
			before: []time.Time{
				newTime(t, "2018-09-01 08:00"),
			},
			after: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-01 12:00"),
			},
			wantErr: false,
		},
		{
			name: "end second day",
			end:  newTime(t, "2018-09-02 16:00"),
			before: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-02 10:00"),
			},
			after: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-02 10:00"),
				newTime(t, "2018-09-02 16:00"),
			},
			wantErr: false,
		},
		{
			name:    "not started",
			end:     newTime(t, "2018-09-01 16:00"),
			before:  []time.Time{},
			after:   []time.Time{},
			wantErr: true,
		},
		{
			name: "new day previous not stopped",
			end:  newTime(t, "2018-09-02 16:00"),
			before: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-02 10:00"),
			},
			after: []time.Time{
				newTime(t, "2018-09-01 08:00"),
				newTime(t, "2018-09-02 10:00"),
				newTime(t, "2018-09-02 16:00"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &timeSheet{times: tt.before}
			err := ts.End(tt.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeSheet.End() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.after, ts.times); diff != "" {
				t.Errorf("timeSheet.End() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_timeSheet_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []time.Time
		wantErr bool
	}{
		{
			name:    "invalid json",
			input:   "testdata/invalid.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid date",
			input:   "testdata/invalid_date.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid time",
			input:   "testdata/invalid_time.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "testdata/empty.json",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty day",
			input:   "testdata/empty_day.json",
			want:    nil,
			wantErr: false,
		},
		{
			name:  "one day only start",
			input: "testdata/one_day_only_start.json",
			want: []time.Time{
				newTime(t, "2018-09-01 10:00"),
			},
			wantErr: false,
		},
		{
			name:  "one day start/end",
			input: "testdata/one_day_start_end.json",
			want: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
			},
			wantErr: false,
		},
		{
			name:  "one day start/end start",
			input: "testdata/one_day_start_end_start.json",
			want: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-01 13:00"),
			},
			wantErr: false,
		},
		{
			name:  "multiple days",
			input: "testdata/multiple_days.json",
			want: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-02 08:00"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := readFile(t, tt.input)

			ts := &timeSheet{}
			err := json.Unmarshal(bytes, &ts)

			if (err != nil) != tt.wantErr {
				t.Errorf("timeSheet.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, ts.times); diff != "" {
				t.Errorf("timeSheet.UnmarshalJSON() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_timeSheet_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		times   []time.Time
		golden  string
		wantErr bool
	}{
		{
			name:    "empty",
			times:   []time.Time{},
			golden:  "testdata/empty.json",
			wantErr: false,
		},
		{
			name:    "one day only start",
			times:   []time.Time{newTime(t, "2018-09-01 10:00")},
			golden:  "testdata/one_day_only_start.json",
			wantErr: false,
		},
		{
			name: "one day start/end",
			times: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
			},
			golden:  "testdata/one_day_start_end.json",
			wantErr: false,
		},
		{
			name: "one day start/end start",
			times: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-01 13:00"),
			},
			golden:  "testdata/one_day_start_end_start.json",
			wantErr: false,
		},
		{
			name: "multiple days",
			times: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-02 08:00"),
			},
			golden:  "testdata/multiple_days.json",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &timeSheet{times: tt.times}

			bytes, err := json.MarshalIndent(ts, "", "  ")

			if (err != nil) != tt.wantErr {
				t.Errorf("timeSheet.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			want := readFile(t, tt.golden)
			if diff := cmp.Diff(want, bytes); diff != "" {
				t.Errorf("timeSheet.MarshalJSON() differs: (-want +got)\n%s", diff)
			}
		})
	}
}
