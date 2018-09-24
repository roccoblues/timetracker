package main

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_newFileRepo(t *testing.T) {
	want := &fileRepo{path: "foo/bar"}
	if got := newFileRepo("foo/bar"); !reflect.DeepEqual(got, want) {
		t.Errorf("newFileRepo() = %v, want %v", got, want)
	}
}

func Test_fileRepo_Read(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    []time.Time
		wantErr bool
	}{
		{
			name:    "non-existing file",
			path:    "non-existing.json",
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "invalid json",
			path:    "testdata/invalid.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid date",
			path:    "testdata/invalid_date.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid time",
			path:    "testdata/invalid_time.json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			path:    "testdata/empty.json",
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "empty day",
			path:    "testdata/empty_day.json",
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name: "one day only start",
			path: "testdata/one_day_only_start.json",
			want: []time.Time{
				newTime(t, "2018-09-01 10:00"),
			},
			wantErr: false,
		},
		{
			name: "one day start/end",
			path: "testdata/one_day_start_end.json",
			want: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
			},
			wantErr: false,
		},
		{
			name: "one day start/end start",
			path: "testdata/one_day_start_end_start.json",
			want: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
				newTime(t, "2018-09-01 13:00"),
			},
			wantErr: false,
		},
		{
			name: "multiple days",
			path: "testdata/multiple_days.json",
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
			f := newFileRepo(tt.path)

			got, err := f.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("fileRepo.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("fileRepo.Read() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func Test_fileRepo_Write(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		times   []time.Time
		golden  string
		wantErr bool
	}{
		{
			name:    "creates new file",
			path:    nonExistingFile(t),
			times:   []time.Time{newTime(t, "2018-09-01 10:00")},
			golden:  "testdata/one_day_only_start.json",
			wantErr: false,
		},
		{
			name:    "overwrite existing file",
			path:    newFile(t, []byte("foo bar")),
			times:   []time.Time{newTime(t, "2018-09-01 10:00")},
			golden:  "testdata/one_day_only_start.json",
			wantErr: false,
		},
		{
			name:    "empty",
			path:    nonExistingFile(t),
			times:   []time.Time{},
			golden:  "testdata/empty.json",
			wantErr: false,
		},
		{
			name: "one day start/end",
			path: nonExistingFile(t),
			times: []time.Time{
				newTime(t, "2018-09-01 10:00"),
				newTime(t, "2018-09-01 12:00"),
			},
			golden:  "testdata/one_day_start_end.json",
			wantErr: false,
		},
		{
			name: "one day start/end start",
			path: nonExistingFile(t),
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
			path: nonExistingFile(t),
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
			defer os.Remove(tt.path)
			f := newFileRepo(tt.path)

			if err := f.Write(tt.times); (err != nil) != tt.wantErr {
				t.Errorf("fileRepo.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := readFile(t, f.path)
			want := readFile(t, tt.golden)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("fileRepo.Write() differs: (-want +got)\n%s", diff)
			}
		})
	}
}
