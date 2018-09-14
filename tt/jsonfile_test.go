package main

import (
	"reflect"
	"testing"
	"time"
)

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
			json:    []byte("{\"2018-09-03\":[]}"),
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "one day only start",
			json:    []byte("{\"2018-09-01\":[\"10:00\"]}"),
			want:    []time.Time{testStartTime},
			wantErr: false,
		},
		{
			name:    "one day with start/end",
			json:    []byte("{\"2018-09-01\":[\"10:00\", \"12:00\"]}"),
			want:    []time.Time{testStartTime, testEndTime},
			wantErr: false,
		},
		{
			name:    "one day with start/end start",
			json:    []byte("{\"2018-09-01\":[\"10:00\", \"12:00\", \"13:00\"]}"),
			want:    []time.Time{testStartTime, testEndTime, testStartTime2},
			wantErr: false,
		},
		{
			name:    "multiple days",
			json:    []byte("{\"2018-09-01\":[\"10:00\", \"12:00\"],\"2018-09-02\":[\"08:00\"]}"),
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
			json:    []byte("{\"2018-09-00\":[\"10:00\"]}"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid time",
			json:    []byte("{\"2018-09-01\":[\"25:00\"]}"),
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
