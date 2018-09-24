package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_writeDays(t *testing.T) {
	testTime := newTime(t, "2018-09-01 10:00")

	days := []*day{
		&day{
			Date: testTime,
			Entries: []*entry{
				&entry{
					Start: testTime,
					End:   testTime.Add(time.Hour * 2),
				},
				&entry{
					Start: testTime.Add(time.Hour * 4),
				},
			},
		},
		&day{
			Date: testTime.Add(time.Hour * 24),
			Entries: []*entry{
				&entry{
					Start: testTime.Add(time.Hour * 24),
				},
			},
		},
		&day{
			Date: testTime.Add(time.Hour * 24 * 8),
			Entries: []*entry{
				&entry{
					Start: testTime.Add(time.Hour * 24 * 8),
					End:   testTime.Add(time.Hour * 25 * 8),
				},
			},
		},
	}

	tests := []struct {
		name string
		days []*day
		want string
	}{
		{
			name: "works",
			days: days,
			want: "01.09.2018  2.00  10:00-12:00 14:00-\n02.09.2018  0.00  10:00-\n\n09.09.2018  8.00  10:00-18:00 \n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			writeDays(tt.days, &buf)
			if diff := cmp.Diff(tt.want, string(buf.Bytes())); diff != "" {
				t.Errorf("writeDays() differs: (-want +got)\n%s", diff)
			}
		})
	}
}
