package timesheet

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_Sheet_Start(t *testing.T) {
	tests := []struct {
		name    string
		start   time.Time
		before  []time.Time
		after   []time.Time
		wantErr bool
	}{
		{
			name:   "first entry",
			start:  time.Date(2018, time.September, 1, 10, 0, 0, 0, time.Now().Location()),
			before: []time.Time{},
			after: []time.Time{
				time.Date(2018, time.September, 1, 10, 0, 0, 0, time.Now().Location()),
			},
			wantErr: false,
		},
		{
			name:  "second entry",
			start: time.Date(2018, time.September, 1, 13, 0, 0, 0, time.Now().Location()),
			before: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
			},
			after: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 13, 0, 0, 0, time.Now().Location()),
			},
			wantErr: false,
		},
		{
			name:  "new day",
			start: time.Date(2018, time.September, 2, 8, 0, 0, 0, time.Now().Location()),
			before: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
			},
			after: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 2, 8, 0, 0, 0, time.Now().Location()),
			},
			wantErr: false,
		},
		{
			name:  "new day previous not stopped",
			start: time.Date(2018, time.September, 2, 8, 0, 0, 0, time.Now().Location()),
			before: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
			},
			after: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 2, 8, 0, 0, 0, time.Now().Location()),
			},
			wantErr: false,
		},
		{
			name:  "already started",
			start: time.Date(2018, time.September, 1, 16, 0, 0, 0, time.Now().Location()),
			before: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
			},
			after: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
			},
			wantErr: true,
		},
		{
			name:  "start earlier as end",
			start: time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
			before: []time.Time{
				time.Date(2018, time.September, 1, 9, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 16, 0, 0, 0, time.Now().Location()),
			},
			after: []time.Time{
				time.Date(2018, time.September, 1, 9, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 16, 0, 0, 0, time.Now().Location()),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sheet := &Sheet{Times: tt.before}
			err := sheet.Start(tt.start)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeSheet.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.after, sheet.Times); diff != "" {
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
			end:  time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
			before: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
			},
			after: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
			},
			wantErr: false,
		},
		{
			name: "end second day",
			end:  time.Date(2018, time.September, 2, 16, 0, 0, 0, time.Now().Location()),
			before: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 2, 10, 0, 0, 0, time.Now().Location()),
			},
			after: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 2, 10, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 2, 16, 0, 0, 0, time.Now().Location()),
			},
			wantErr: false,
		},
		{
			name:    "not started",
			end:     time.Date(2018, time.September, 1, 16, 0, 0, 0, time.Now().Location()),
			before:  []time.Time{},
			after:   []time.Time{},
			wantErr: true,
		},
		{
			name: "new day previous not stopped",
			end:  time.Date(2018, time.September, 2, 16, 0, 0, 0, time.Now().Location()),
			before: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 2, 10, 0, 0, 0, time.Now().Location()),
			},
			after: []time.Time{
				time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 2, 10, 0, 0, 0, time.Now().Location()),
				time.Date(2018, time.September, 2, 16, 0, 0, 0, time.Now().Location()),
			},
			wantErr: false,
		},
		{
			name: "end ealier as start",
			end:  time.Date(2018, time.September, 1, 8, 0, 0, 0, time.Now().Location()),
			before: []time.Time{
				time.Date(2018, time.September, 1, 9, 0, 0, 0, time.Now().Location()),
			},
			after: []time.Time{
				time.Date(2018, time.September, 1, 9, 0, 0, 0, time.Now().Location()),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sheet := &Sheet{Times: tt.before}
			err := sheet.End(tt.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeSheet.End() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.after, sheet.Times); diff != "" {
				t.Errorf("timeSheet.End() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestSameDate(t *testing.T) {
	testTime := time.Date(2018, time.September, 1, 10, 0, 0, 0, time.Now().Location())

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
