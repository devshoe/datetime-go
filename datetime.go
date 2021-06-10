package datetime

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

// ParseDatetime takes a datetime string, and returns a time.Time
// The input must only follow year-month-day hour:min:sec decendance.
// Valid inputs look like this:
//"2020-12-12", "2020-12-12T20:20:18Z", "1999-1-7 09:16:28", "1779-12-23 09:15", "20-12-21", "2020-12-12T20:20:18"
func ParseDatetime(datetime string) (time.Time, error) {
	return dateparse.ParseAny(datetime)
}

//InlineParseDatetime is ParseDatetime without the error, to allow for more succint code. It fails silently
func InlineParseDatetime(datetime string) time.Time {
	t, _ := ParseDatetime(datetime)
	return t
}

//ParseDatetimeArray parses a row of datetime strings
func ParseDatetimeArray(datetimes []string) ([]time.Time, error) {
	parsed := []time.Time{}
	for _, d := range datetimes {
		t, err := ParseDatetime(d)
		if err != nil {
			return parsed, err
		}
		parsed = append(parsed, t)
	}
	return parsed, nil
}

//InlineParseDatetimeArray is ParseDatetimeArray without the error, to allow for more succint code. It fails silently
func InlineParseDatetimeArray(datetimes []string) []time.Time {
	parsed := []time.Time{}
	for _, d := range datetimes {
		t, _ := ParseDatetime(d)

		parsed = append(parsed, t)
	}
	return parsed
}

//ParseDatetimeWithYYMMDDLikeLayout uses a layout of YYMMDD type rather than numeric in golang time lib
func ParseDatetimeWithYYMMDDLikeLayout(datetime string, layout string) (time.Time, error) {
	return time.Parse(convertLayoutToGolangLayout(layout), datetime)
}

var matchIntervalWithDigit, _ = regexp.Compile("[0-9]+[a-z]{1}")
var matchIntervalWithoutDigit, _ = regexp.Compile("^[a-z]")

//ParseInterval parses an interval string like "1minute", "minute", "1m"
func ParseInterval(interval string) (time.Duration, error) {
	match := matchIntervalWithoutDigit.FindString(interval)
	if match != "" {
		interval = "1" + interval
	}
	match = matchIntervalWithDigit.FindString(interval)
	if match == "" {
		return 0, fmt.Errorf("(ParseInterval) failed for %v: could not find match", interval)
	}
	switch match[len(match)-1] {
	case 'd':
		i, _ := strconv.Atoi(match[:len(match)-1])
		match = strconv.Itoa(i*24) + "h"
	case 'w':
		i, _ := strconv.Atoi(match[:len(match)-1])
		match = strconv.Itoa(i*24*7) + "h"
	case 'y':
		i, _ := strconv.Atoi(match[:len(match)-1])
		match = strconv.Itoa(i*24*7*365) + "h"
	}
	return time.ParseDuration(match)
}

//InlineParseInterval is ParseInterval without error. Fails silently
func InlineParseInterval(interval string) time.Duration {
	d, _ := ParseInterval(interval)
	return d
}

//ParseAbstract parses the following words: "today", "yesterday", "1week", "2week", "3week", "1month", "2month", "1year"
func ParseAbstract(absString string) (time.Time, error) {
	var parsed time.Time
	switch absString {
	case "today":
		parsed = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	case "yesterday":
		parsed = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.Now().Location())
	case "1week":
		parsed = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-7, 0, 0, 0, 0, time.Now().Location())
	case "1month":
		parsed = time.Date(time.Now().Year(), time.Now().Month()-1, time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	case "2month":
		parsed = time.Date(time.Now().Year(), time.Now().Month()-2, time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	case "1year":
		parsed = time.Date(time.Now().Year()-1, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	default:
		return time.Time{}, fmt.Errorf("(ParseAbstract)couldnt parse: %v, invalid string", absString)
	}
	return parsed, nil
}

//ParseStringOrTime takes a date of either string or time.Time and converts them to an array of time.Time
//All dates are stripped of location data, but no arithmetic is performed
//Basically, a date like "2021-12-12 09:15:00+0530" becomes "2021-12-12 09:15:00+0000"
func ParseStringOrTime(date interface{}) (time.Time, error) {
	var parsed time.Time
	if date == nil {
		return time.Time{}, fmt.Errorf("(ParseDate) nil interface  >.<")
	}
	switch date := date.(type) {
	case string:
		t, err := dateparse.ParseAny(date)
		if err != nil {
			return time.Time{}, err
		}
		parsed = t
	case time.Time:
		parsed = date
	default:
		return time.Time{}, fmt.Errorf("(ParseDate) parsing %v failed, invalid type", date)
	}

	if parsed.Location() != time.UTC {
		parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), parsed.Nanosecond(), time.UTC)
	}

	return parsed, nil
}
