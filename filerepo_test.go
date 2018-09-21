package main

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/roccoblues/tt/test"
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
			path:    test.NonExistingFile(t),
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "invalid json",
			path:    test.NewFile(t, []byte(test.InvalidJSON)),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid date",
			path:    test.NewFile(t, []byte(test.InvalidDateJSON)),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid time",
			path:    test.NewFile(t, []byte(test.InvalidTimeJSON)),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty",
			path:    test.NewFile(t, []byte(test.EmptyJSON)),
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "empty day",
			path:    test.NewFile(t, []byte(test.EmptyDayJSON)),
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name: "one day only start",
			path: test.NewFile(t, []byte(test.OneDayOnlyStartJSON)),
			want: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
			},
			wantErr: false,
		},
		{
			name: "one day start/end",
			path: test.NewFile(t, []byte(test.OneDayStartEndJSON)),
			want: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
				test.Time(t, "2018-09-01 12:00"),
			},
			wantErr: false,
		},
		{
			name: "one day start/end start",
			path: test.NewFile(t, []byte(test.OneDayStartEndStartJSON)),
			want: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-01 13:00"),
			},
			wantErr: false,
		},
		{
			name: "multiple days",
			path: test.NewFile(t, []byte(test.MultipleDaysJSON)),
			want: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-02 08:00"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.path)
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
		want    []byte
		wantErr bool
	}{
		{
			name:    "creates new file",
			path:    test.NonExistingFile(t),
			times:   []time.Time{test.Time(t, "2018-09-01 10:00")},
			want:    []byte(test.OneDayOnlyStartJSON),
			wantErr: false,
		},
		{
			name:    "overwrite existing file",
			path:    test.NewFile(t, []byte("foo bar")),
			times:   []time.Time{test.Time(t, "2018-09-01 10:00")},
			want:    []byte(test.OneDayOnlyStartJSON),
			wantErr: false,
		},
		{
			name:    "empty",
			path:    test.NonExistingFile(t),
			times:   []time.Time{},
			want:    []byte(test.EmptyJSON),
			wantErr: false,
		},
		{
			name: "one day only start",
			path: test.NonExistingFile(t),
			times: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
			},
			want:    []byte(test.OneDayOnlyStartJSON),
			wantErr: false,
		},
		{
			name: "one day start/end",
			path: test.NonExistingFile(t),
			times: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
				test.Time(t, "2018-09-01 12:00"),
			},
			want:    []byte(test.OneDayStartEndJSON),
			wantErr: false,
		},
		{
			name: "one day start/end start",
			path: test.NonExistingFile(t),
			times: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-01 13:00"),
			},
			want:    []byte(test.OneDayStartEndStartJSON),
			wantErr: false,
		},
		{
			name: "multiple days",
			path: test.NonExistingFile(t),
			times: []time.Time{
				test.Time(t, "2018-09-01 10:00"),
				test.Time(t, "2018-09-01 12:00"),
				test.Time(t, "2018-09-02 08:00"),
			},
			want:    []byte(test.MultipleDaysJSON),
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
			got := test.ReadFile(t, f.path)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("fileRepo.Write() differs: (-want +got)\n%s", diff)
			}
		})
	}
}
