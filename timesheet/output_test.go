package timesheet

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

var outputTestCases = []struct {
	description string
	times       []time.Time
	fixture     string
}{
	{
		description: "empty",
		times:       []time.Time{},
		fixture:     "test-fixtures/output_empty.txt",
	},
	{
		description: "august",
		times: []time.Time{
			time.Date(2018, time.August, 28, 8, 0, 0, 0, time.Now().Location()),
			time.Date(2018, time.August, 28, 12, 0, 0, 0, time.Now().Location()),
		},
		fixture: "test-fixtures/output_august.txt",
	},
	{
		description: "september",
		times: []time.Time{
			time.Date(2018, time.September, 1, 10, 0, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 1, 11, 42, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 1, 14, 0, 0, 0, time.Now().Location()),

			time.Date(2018, time.September, 2, 8, 0, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 2, 16, 0, 0, 0, time.Now().Location()),

			time.Date(2018, time.September, 9, 8, 0, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 9, 12, 24, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 9, 13, 12, 0, 0, time.Now().Location()),
			time.Date(2018, time.September, 9, 17, 57, 0, 0, time.Now().Location()),
		},
		fixture: "test-fixtures/output_september.txt",
	},
}

func TestPrint(t *testing.T) {
	var timeFormat = "15:04"
	var dateFormat = "02.01.2006"

	for _, tc := range outputTestCases {
		t.Run(tc.description, func(t *testing.T) {
			output := &bytes.Buffer{}

			print(tc.times, 15*time.Minute, dateFormat, timeFormat, output)

			want := string(readFile(t, tc.fixture))
			if diff := cmp.Diff(want, output.String()); diff != "" {
				t.Errorf("Print() differs: (-want +got)\n%s", diff)
			}
		})
	}
}
