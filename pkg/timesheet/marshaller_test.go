package timesheet

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var marshalTestCases = []struct {
	description string
	fixture     string
	times       []time.Time
	wantErr     bool
	skipMarshal bool
}{
	{
		description: "invalid json",
		fixture:     "testdata/invalid.json",
		times:       nil,
		wantErr:     true,
		skipMarshal: true,
	},
	{
		description: "invalid date",
		fixture:     "testdata/invalid_date.json",
		times:       nil,
		wantErr:     true,
		skipMarshal: true,
	},
	{
		description: "invalid time",
		fixture:     "testdata/invalid_time.json",
		times:       nil,
		wantErr:     true,
		skipMarshal: true,
	},
	{
		description: "empty file",
		fixture:     "testdata/empty.json",
		skipMarshal: true,
		times:       nil,
	},
	{
		description: "empty json",
		fixture:     "testdata/empty_json.json",
		times:       nil,
	},
	{
		description: "empty day",
		fixture:     "testdata/empty_day.json",
		times:       nil,
		skipMarshal: true,
	},
	{
		description: "one day only start",
		fixture:     "testdata/one_day_only_start.json",
		times: []time.Time{
			time.Date(2018, time.September, 1, 10, 0, 0, 0, time.Now().Location()),
		},
	},
	{
		description: "one day start/end",
		fixture:     "testdata/one_day_start_end.json",
		times: []time.Time{
			time.Date(2018, time.September, 1, 10, 0, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
		},
	},
	{
		description: "one day start/end start",
		fixture:     "testdata/one_day_start_end_start.json",
		times: []time.Time{
			time.Date(2018, time.September, 1, 10, 0, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 1, 13, 0, 0, 0, time.Now().Location()),
		},
	},
	{
		description: "multiple days",
		fixture:     "testdata/multiple_days.json",
		times: []time.Time{
			time.Date(2018, time.September, 1, 10, 0, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 1, 12, 0, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 2, 8, 0, 0, 0, time.Now().Location()),
		},
	},
}

func TestUnmarshal(t *testing.T) {
	var timeFormat = "15:04"
	var dateFormat = "02.01.2006"

	for _, tc := range marshalTestCases {
		t.Run(tc.description, func(t *testing.T) {
			file, _ := os.Open(tc.fixture)

			actual, err := unmarshal(file, dateFormat, timeFormat)

			if (err != nil) != tc.wantErr {
				t.Errorf("unmarshal() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.times, actual); diff != "" {
				t.Errorf("unmarshal() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	var timeFormat = "15:04"
	var dateFormat = "02.01.2006"

	for _, tc := range marshalTestCases {
		if tc.skipMarshal {
			continue
		}

		t.Run(tc.description, func(t *testing.T) {
			want := readFile(t, tc.fixture)

			var actual bytes.Buffer
			marshal(&actual, tc.times, dateFormat, timeFormat)

			if diff := cmp.Diff(strings.Replace(string(want), "\r\n", "\n", -1), strings.Replace(actual.String(), "\r\n", "\n", -1)); diff != "" {
				t.Errorf("marshal() differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	var timeFormat = "15:04"
	var dateFormat = "02.01.2006"

	for _, tc := range marshalTestCases {
		if tc.skipMarshal {
			continue
		}

		t.Run(tc.description, func(t *testing.T) {
			var actual bytes.Buffer

			marshal(&actual, tc.times, dateFormat, timeFormat)
			times, err := unmarshal(&actual, dateFormat, timeFormat)

			if err != nil {
				t.Errorf("unmarshal(marshal()) error = %v", err)
				return
			}
			if diff := cmp.Diff(tc.times, times); diff != "" {
				t.Errorf("unmarshal(marshal()) differs: (-want +got)\n%s", diff)
			}
		})
	}

}
