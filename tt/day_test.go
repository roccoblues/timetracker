package main

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_day_addStartTime(t *testing.T) {
	testTime, err := time.Parse("2006-01-02 15:04", "2018-09-11 13:40")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		day     day
		time    time.Time
		wantErr bool
		times   []time.Time
	}{
		{
			name:    "empty",
			day:     day{},
			time:    testTime,
			wantErr: false,
			times:   []time.Time{testTime},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.day.addStartTime(tt.time); (err != nil) != tt.wantErr {
				t.Errorf("day.addStartTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := cmp.Diff(tt.day.times, tt.times); diff != "" {
				t.Errorf("day.addStartTime() differ: (-want +got)\n%s", diff)
			}
		})
	}
}
