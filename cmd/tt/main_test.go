package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	now := time.Now()
	timeFormat := "15:04"
	dateFormat := "02.01.2006"

	tests := []struct {
		name    string
		value   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "works",
			value:   "10:15",
			want:    time.Date(now.Year(), now.Month(), now.Day(), 10, 15, 0, 0, now.Location()),
			wantErr: false,
		},
		{
			name:    "fails",
			value:   "99:15",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTime(tt.value, dateFormat, timeFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
