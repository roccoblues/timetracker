package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
)

var emptyDay = heredoc.Doc(`{
  "2018-09-01": []
}`)
var oneDayOnlyStart = heredoc.Doc(`{
  "2018-09-01": [
    "10:00"
  ]
}`)
var oneDayStartEnd = heredoc.Doc(`{
  "2018-09-01": [
    "10:00",
    "12:00"
  ]
}`)
var oneDayStartEndStart = heredoc.Doc(`{
  "2018-09-01": [
    "10:00",
    "12:00",
    "13:00"
  ]
}`)
var multipleDays = heredoc.Doc(`{
  "2018-09-01": [
    "10:00",
    "12:00"
  ],
  "2018-09-02": [
    "08:00"
  ]
}`)
var invalidDate = heredoc.Doc(`{
  "2018-09-00": [
    "10:00"
  ]
}`)
var invalidTime = heredoc.Doc(`{
  "2018-09-01": [
    "25:00"
  ]
}`)

func Test_decode(t *testing.T) {
	testStartTime, err := time.Parse("2006-01-02 15:04", "2018-09-01 10:00")
	if err != nil {
		t.Fatal(err)
	}
	testEndTime, err := time.Parse("2006-01-02 15:04", "2018-09-01 12:00")
	if err != nil {
		t.Fatal(err)
	}
	testStartTime2, err := time.Parse("2006-01-02 15:04", "2018-09-01 13:00")
	if err != nil {
		t.Fatal(err)
	}
	testStartTime3, err := time.Parse("2006-01-02 15:04", "2018-09-02 08:00")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		json    []byte
		want    []time.Time
		wantErr bool
	}{
		{
			name:    "empty",
			json:    []byte("{}"),
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "empty day",
			json:    []byte(emptyDay),
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "one day only start",
			json:    []byte(oneDayOnlyStart),
			want:    []time.Time{testStartTime},
			wantErr: false,
		},
		{
			name:    "one day with start/end",
			json:    []byte(oneDayStartEnd),
			want:    []time.Time{testStartTime, testEndTime},
			wantErr: false,
		},
		{
			name:    "one day with start/end start",
			json:    []byte(oneDayStartEndStart),
			want:    []time.Time{testStartTime, testEndTime, testStartTime2},
			wantErr: false,
		},
		{
			name:    "multiple days",
			json:    []byte(multipleDays),
			want:    []time.Time{testStartTime, testEndTime, testStartTime3},
			wantErr: false,
		},
		{
			name:    "invalid json",
			json:    []byte("{"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid date",
			json:    []byte(invalidDate),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid time",
			json:    []byte(invalidTime),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decode(tt.json)
			if (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encode(t *testing.T) {
	testStartTime, err := time.Parse("2006-01-02 15:04", "2018-09-01 10:00")
	if err != nil {
		t.Fatal(err)
	}
	testEndTime, err := time.Parse("2006-01-02 15:04", "2018-09-01 12:00")
	if err != nil {
		t.Fatal(err)
	}
	testStartTime2, err := time.Parse("2006-01-02 15:04", "2018-09-01 13:00")
	if err != nil {
		t.Fatal(err)
	}
	testStartTime3, err := time.Parse("2006-01-02 15:04", "2018-09-02 08:00")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		times   []time.Time
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty",
			times:   []time.Time{},
			want:    []byte("{}"),
			wantErr: false,
		},
		{
			name:    "one day only start",
			times:   []time.Time{testStartTime},
			want:    []byte(oneDayOnlyStart),
			wantErr: false,
		},
		{
			name:    "one day with start/end",
			times:   []time.Time{testStartTime, testEndTime},
			want:    []byte(oneDayStartEnd),
			wantErr: false,
		},
		{
			name:    "one day with start/end start",
			times:   []time.Time{testStartTime, testEndTime, testStartTime2},
			want:    []byte(oneDayStartEndStart),
			wantErr: false,
		},
		{
			name:    "multiple days",
			times:   []time.Time{testStartTime, testEndTime, testStartTime3},
			want:    []byte(multipleDays),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encode(tt.times)
			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encode() = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_encode_decode(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "one day only start",
			json: oneDayOnlyStart,
		},
		{
			name: "one day with start/end",
			json: oneDayStartEnd,
		},
		{
			name: "one day with start/end start",
			json: oneDayStartEndStart,
		},
		{
			name: "multiple days",
			json: multipleDays,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := decode([]byte(tt.json))
			if err != nil {
				t.Fatal(err)
			}
			got, err := encode(decoded)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.json {
				t.Errorf("encode(decode(%s)) = %s", tt.json, got)
			}
		})
	}
}

func Test_jsonFile_Read(t *testing.T) {
	dir, err := ioutil.TempDir("", "jsonfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	nonExistingFile := filepath.Join(dir, "test")

	tests := []struct {
		name    string
		path    string
		want    []time.Time
		wantErr bool
	}{
		{
			name:    "file doesn't exist",
			path:    nonExistingFile,
			want:    []time.Time{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := newJSONFile(tt.path)
			got, err := j.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonFile.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonFile.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
