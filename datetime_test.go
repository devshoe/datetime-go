package datetime

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var validSamples = []string{"2020-12-12", "2020-12-12T20:20:18Z", "1999-1-7 09:16:28", "1779-12-23 09:15", "20-21-12", "2020-12-12T20:20:18"}
var validRFC3339Equivalents = []string{"2020-12-12T00:00:00Z", "2020-12-12T20:20:18Z", "1999-01-07T09:16:28Z", "1779-12-23T09:15:00Z", "2020-12-21T00:00:00Z", "2020-12-12T20:20:18Z"}
var invalidSamples = []string{"2020-12-", "20 12-12", "123", "1-12-3 23:2-2"}

// //RFC3339NoTZ is like 2021-12-12T09:05:12Z, basically nothing after Z
const RFC3339NoTZ = "2006-01-02T15:04:05Z"

//input rfc string, return without error
func parseDatetimeNoError(dt string) time.Time {
	d, _ := time.Parse(RFC3339NoTZ, dt)
	return d
}

func TestParseDatetime(t *testing.T) {
	type testcase struct {
		inDatetime   string
		gotDatetime  time.Time
		gotErr       error
		wantDatetime time.Time
		wantErr      error
	}
	cases := []testcase{}
	for i, s := range validSamples {
		gotDatetime, gotErr := ParseDatetime(s)
		cases = append(cases, testcase{s, gotDatetime, gotErr, parseDatetimeNoError(validRFC3339Equivalents[i]), nil})
	}
	for _, s := range invalidSamples {
		gotDatetime, gotErr := ParseDatetime(s)
		cases = append(cases, testcase{s, gotDatetime, gotErr, time.Time{}, fmt.Errorf("(ParseDatetime)invalid datetime, nothing found for %v", s)})
	}
	for _, c := range cases {
		if c.gotDatetime != c.wantDatetime {
			t.Errorf("(TestParseDatetime) test failed, wanted output datetime:%v, got: %v", c.wantDatetime, c.gotDatetime)
		}
		if c.gotErr != c.wantErr {
			t.Errorf("(TestParseDatetime) test failed, wanted output error %v, got %v", c.wantErr, c.gotErr)
		}
	}
}

func TestParseDatetimeWithLayout(t *testing.T) {
	type args struct {
		datetime string
		layout   string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{"only date", args{"2020-12-12", "YYYY-MM-DD"}, parseDatetimeNoError(validRFC3339Equivalents[0]), false},
		{"rfc3339", args{"2020-12-12T20:20:18Z", "YYYY-MM-DDThh:mm:ssZ"}, parseDatetimeNoError(validRFC3339Equivalents[1]), false},
		{"date space time", args{"1999-01-07 09:16:28", "YYYY-MM-DD hh:mm:ss"}, parseDatetimeNoError(validRFC3339Equivalents[2]), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDatetimeWithYYMMDDLikeLayout(tt.args.datetime, tt.args.layout)
			t.Log(tt.args.datetime, tt.args.layout, got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDatetimeWithLayout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDatetimeWithLayout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInterval(t *testing.T) {
	type args struct {
		interval string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{"withdigits", args{"1minute"}, time.Minute, false},
		{"withoutdigits", args{"minute"}, time.Minute, false},
		{"hours", args{"1hour"}, time.Hour, false},
		{"daynodigit", args{"da"}, time.Hour * 24, false},
		{"day", args{"1day"}, time.Hour * 24, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInterval(tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTimeFromDatetime(t *testing.T) {
	type args struct {
		datetime time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractTimeFromDatetime(tt.args.datetime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractTimeFromDatetime() = %v, want %v", got, tt.want)
			}
		})
	}
}
