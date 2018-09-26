package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
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

func Test_sameDate(t *testing.T) {
	testTime := newTime(t, "01.09.2018 10:00")

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
			start:  newTime(t, "01.09.2018 10:00"),
			before: []time.Time{},
			after: []time.Time{
				newTime(t, "01.09.2018 10:00"),
			},
			wantErr: false,
		},
		{
			name:  "second entry",
			start: newTime(t, "01.09.2018 13:00"),
			before: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "01.09.2018 12:00"),
			},
			after: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "01.09.2018 12:00"),
				newTime(t, "01.09.2018 13:00"),
			},
			wantErr: false,
		},
		{
			name:  "new day",
			start: newTime(t, "02.09.2018 08:00"),
			before: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "01.09.2018 12:00"),
			},
			after: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "01.09.2018 12:00"),
				newTime(t, "02.09.2018 08:00"),
			},
			wantErr: false,
		},
		{
			name:  "new day previous not stopped",
			start: newTime(t, "02.09.2018 08:00"),
			before: []time.Time{
				newTime(t, "01.09.2018 08:00"),
			},
			after: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "02.09.2018 08:00"),
			},
			wantErr: false,
		},
		{
			name:  "already started",
			start: newTime(t, "01.09.2018 16:00"),
			before: []time.Time{
				newTime(t, "01.09.2018 08:00"),
			},
			after: []time.Time{
				newTime(t, "01.09.2018 08:00"),
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
			end:  newTime(t, "01.09.2018 12:00"),
			before: []time.Time{
				newTime(t, "01.09.2018 08:00"),
			},
			after: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "01.09.2018 12:00"),
			},
			wantErr: false,
		},
		{
			name: "end second day",
			end:  newTime(t, "02.09.2018 16:00"),
			before: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "01.09.2018 12:00"),
				newTime(t, "02.09.2018 10:00"),
			},
			after: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "01.09.2018 12:00"),
				newTime(t, "02.09.2018 10:00"),
				newTime(t, "02.09.2018 16:00"),
			},
			wantErr: false,
		},
		{
			name:    "not started",
			end:     newTime(t, "01.09.2018 16:00"),
			before:  []time.Time{},
			after:   []time.Time{},
			wantErr: true,
		},
		{
			name: "new day previous not stopped",
			end:  newTime(t, "02.09.2018 16:00"),
			before: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "02.09.2018 10:00"),
			},
			after: []time.Time{
				newTime(t, "01.09.2018 08:00"),
				newTime(t, "02.09.2018 10:00"),
				newTime(t, "02.09.2018 16:00"),
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
				newTime(t, "01.09.2018 10:00"),
			},
			wantErr: false,
		},
		{
			name:  "one day start/end",
			input: "testdata/one_day_start_end.json",
			want: []time.Time{
				newTime(t, "01.09.2018 10:00"),
				newTime(t, "01.09.2018 12:00"),
			},
			wantErr: false,
		},
		{
			name:  "one day start/end start",
			input: "testdata/one_day_start_end_start.json",
			want: []time.Time{
				newTime(t, "01.09.2018 10:00"),
				newTime(t, "01.09.2018 12:00"),
				newTime(t, "01.09.2018 13:00"),
			},
			wantErr: false,
		},
		{
			name:  "multiple days",
			input: "testdata/multiple_days.json",
			want: []time.Time{
				newTime(t, "01.09.2018 10:00"),
				newTime(t, "01.09.2018 12:00"),
				newTime(t, "02.09.2018 08:00"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := readFile(t, tt.input)

			var ts timeSheet
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
			times:   []time.Time{newTime(t, "01.09.2018 10:00")},
			golden:  "testdata/one_day_only_start.json",
			wantErr: false,
		},
		{
			name: "one day start/end",
			times: []time.Time{
				newTime(t, "01.09.2018 10:00"),
				newTime(t, "01.09.2018 12:00"),
			},
			golden:  "testdata/one_day_start_end.json",
			wantErr: false,
		},
		{
			name: "one day start/end start",
			times: []time.Time{
				newTime(t, "01.09.2018 10:00"),
				newTime(t, "01.09.2018 12:00"),
				newTime(t, "01.09.2018 13:00"),
			},
			golden:  "testdata/one_day_start_end_start.json",
			wantErr: false,
		},
		{
			name: "multiple days",
			times: []time.Time{
				newTime(t, "01.09.2018 10:00"),
				newTime(t, "01.09.2018 12:00"),
				newTime(t, "02.09.2018 08:00"),
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

func Test_loadTimeSheet(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name    string
		path    string
		want    *timeSheet
		wantErr bool
	}{
		{
			name:    "initializes empty timesheet if file doesn't exist",
			path:    filepath.Join(tempDir, "foobar"),
			want:    &timeSheet{},
			wantErr: false,
		},
		{
			name:    "returns an error if json is invalid",
			path:    "testdata/invalid.json",
			want:    nil,
			wantErr: true,
		},
		{
			name: "loads timesheet from json",
			path: "testdata/one_day_only_start.json",
			want: &timeSheet{
				times: []time.Time{
					newTime(t, "01.09.2018 10:00"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadTimeSheet(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadTimeSheet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadTimeSheet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeSheet_Save(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name     string
		ts       *timeSheet
		path     string
		contents string
		wantErr  bool
	}{
		{
			name:     "creates new file",
			ts:       &timeSheet{},
			path:     filepath.Join(tempDir, "foobar"),
			contents: "testdata/empty.json",
			wantErr:  false,
		},
		{
			name: "overwrites existing file",
			ts: &timeSheet{
				times: []time.Time{
					newTime(t, "01.09.2018 10:00"),
				},
			},
			path:     filepath.Join(tempDir, "foobar"),
			contents: "testdata/one_day_only_start.json",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ts.Save(tt.path); (err != nil) != tt.wantErr {
				t.Errorf("timeSheet.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			want := readFile(t, tt.contents)
			got := readFile(t, tt.path)
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("timeSheet.End() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_timeSheet_Print(t *testing.T) {
	times := []time.Time{
		newTime(t, "01.09.2018 10:00"),
		newTime(t, "01.09.2018 11:42"),
		newTime(t, "01.09.2018 14:00"),

		newTime(t, "02.09.2018 08:00"),
		newTime(t, "02.09.2018 16:00"),

		newTime(t, "09.09.2018 08:00"),
		newTime(t, "09.09.2018 12:24"),
		newTime(t, "09.09.2018 13:12"),
		newTime(t, "09.09.2018 17:57"),
	}

	tests := []struct {
		name    string
		ts      *timeSheet
		roundTo time.Duration
		want    string
	}{
		{
			name:    "default",
			ts:      &timeSheet{times: times},
			roundTo: defaultRoundToMinutes * time.Minute,
			want:    "testdata/output_default.txt",
		},
		{
			name:    "round to half hours",
			ts:      &timeSheet{times: times},
			roundTo: 30 * time.Minute,
			want:    "testdata/output_round_to_30_minutes.txt",
		},
		{
			name:    "not rounded",
			ts:      &timeSheet{times: times},
			roundTo: 0 * time.Minute,
			want:    "testdata/output_not_rounded.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			tt.ts.Print(output, tt.roundTo)
			want := string(readFile(t, tt.want))
			if diff := cmp.Diff(output.String(), want); diff != "" {
				t.Errorf("timeSheet.Print() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

// func Test_writeDays(t *testing.T) {
// 	testTime := newTime(t, "01.09.2018 10:00")

// 	days := []*day{
// 		&day{
// 			Date: testTime,
// 			Entries: []*entry{
// 				&entry{
// 					Start: testTime,
// 					End:   testTime.Add(time.Minute * 100),
// 				},
// 				&entry{
// 					Start: testTime.Add(time.Hour * 4),
// 				},
// 			},
// 		},
// 		&day{
// 			Date: testTime.Add(time.Hour * 24),
// 			Entries: []*entry{
// 				&entry{
// 					Start: testTime.Add(time.Hour * 24),
// 				},
// 			},
// 		},
// 		&day{
// 			Date: testTime.Add(time.Hour * 24 * 8),
// 			Entries: []*entry{
// 				&entry{
// 					Start: testTime.Add(time.Hour * 24 * 8),
// 					End:   testTime.Add(time.Hour * 25 * 8),
// 				},
// 			},
// 		},
// 	}

// 	tests := []struct {
// 		name       string
// 		days       []*day
// 		roundTo    time.Duration
// 		dateFormat string
// 		timeFormat string
// 		want       string
// 	}{
// 		{
// 			name:       "default",
// 			days:       days,
// 			want:       "01.09.2018  1.75  10:00-11:45 14:00-\n02.09.2018  0.00  10:00-\n\n09.09.2018  8.00  10:00-18:00 \n",
// 			roundTo:    defaultRoundToMinutes * time.Minute,
// 			dateFormat: dateFormat,
// 			timeFormat: timeFormat,
// 		},
// 		{
// 			name:       "round to half hours",
// 			days:       days,
// 			want:       "01.09.2018  1.50  10:00-11:30 14:00-\n02.09.2018  0.00  10:00-\n\n09.09.2018  8.00  10:00-18:00 \n",
// 			roundTo:    30 * time.Minute,
// 			dateFormat: dateFormat,
// 			timeFormat: timeFormat,
// 		},
// 		{
// 			name:       "different date format",
// 			days:       days,
// 			want:       "2018-01-09  1.75  10:00-11:45 14:00-\n2018-02-09  0.00  10:00-\n\n2018-09-09  8.00  10:00-18:00 \n",
// 			roundTo:    defaultRoundToMinutes * time.Minute,
// 			dateFormat: "2006-02-01",
// 			timeFormat: "15:04",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			buf := bytes.Buffer{}
// 			writeDays(tt.days, &buf, tt.roundTo, tt.dateFormat, tt.timeFormat)
// 			if diff := cmp.Diff(tt.want, string(buf.Bytes())); diff != "" {
// 				t.Errorf("writeDays() differs: (-want +got)\n%s", diff)
// 			}
// 		})
// 	}
// }
