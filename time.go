package datetime

import (
	"time"

	"github.com/sirupsen/logrus"
)

//ExtractTimeFromDatetime returns a time.Time with all the date fields removed ie set to 0-0-0 HH:MM:SS
func ExtractTimeFromDatetime(datetime time.Time) time.Time {
	return time.Date(0, time.January, 0, datetime.Hour(), datetime.Minute(), datetime.Second(), datetime.Nanosecond(), datetime.Location())
}

//ExtractDateFromDatetime sets time of a time.Time to 0 0 0 thus returning yyyy-mm-dd 0:0:0
func ExtractDateFromDatetime(datetime time.Time) time.Time {
	return time.Date(datetime.Year(), datetime.Month(), datetime.Day(), 0, 0, 0, 0, datetime.Location())
}

//ReplaceYear replaces only year with provided integer. If `with` is invalid, it will return original time
func ReplaceYear(t time.Time, with int) time.Time {
	if with > 12 || with < 0 {
		logrus.Errorln("invalid replacement for month: must be between 0 and 12")
		return t
	}
	return time.Date(t.Year(), time.Month(with), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

//ReplaceMonth replaces only month with provided integer. If `with` is invalid, it will return original time
func ReplaceMonth(t time.Time, with int) time.Time {
	if with > 12 || with < 0 {
		logrus.Errorln("invalid replacement for month: must be between 0 and 12")
		return t
	}
	return time.Date(t.Year(), time.Month(with), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

//ReplaceDay replaces only day with provided integer. If `with` is invalid, it will return original time
func ReplaceDay(t time.Time, with int) time.Time {
	if with > 31 || with < 1 {
		logrus.Errorln("invalid replacement for day: must be between 1 and 31")
		return t
	}
	return time.Date(t.Year(), t.Month(), with, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

//ReplaceHour replaces only hour with provided integer. If `with` is invalid, it will return original time
func ReplaceHour(t time.Time, with int) time.Time {
	if with > 23 || with < 0 {
		logrus.Errorln("invalid replacement for hour: must be between 0 and 23")
		return t
	}
	return time.Date(t.Year(), t.Month(), t.Day(), with, t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

//ReplaceMinute replaces only minute with provided integer. If `with` is invalid, it will return original time
func ReplaceMinute(t time.Time, with int) time.Time {
	if with > 59 || with < 0 {
		logrus.Errorln("invalid replacement for minute: must be between 0 and 59")
		return t
	}
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), with, t.Second(), t.Nanosecond(), t.Location())
}

//ReplaceSecond replaces only second with provided integer. If `with` is invalid, it will return original time
func ReplaceSecond(t time.Time, with int) time.Time {
	if with > 59 || with < 0 {
		logrus.Errorln("invalid replacement for second: must be between 0 and 59")
		return t
	}
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), with, t.Nanosecond(), t.Location())
}

//ReplaceNanosecond replaces only nsec with provided integer. If `with` is invalid, it will return original time
func ReplaceNanosecond(t time.Time, with int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), with, t.Location())
}

//ReplaceTimeInDatetime replaces time in date with hour min sec. Nothing else is changed
//if nsec is not provided, it is set to times
func ReplaceTimeInDatetime(t time.Time, hour int, min int, sec int, nsec ...int) time.Time {
	nano := t.Nanosecond()
	if nsec != nil {
		nano = nsec[0]
	}
	return time.Date(t.Year(), t.Month(), t.Day(), hour, min, sec, nano, t.Location())
}

//ReplaceDateInDatetime replaces date in datetime with year month day. Nothing else is changed
func ReplaceDateInDatetime(t time.Time, year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

//StripTimezone removes timezone without adjusting date
//use carefully
func StripTimezone(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}

//DateIsEqual checks for year month and day equality
func DateIsEqual(t1, t2 time.Time) bool {
	return (t1.Year() == t2.Year()) && (t1.Month() == t2.Month()) && (t1.Day() == t2.Day())
}

//TimeIsEqual checks for hour, min, sec equality
//if with nsec is passed, will also compare nano secs
func TimeIsEqual(t1, t2 time.Time, withNsec ...bool) bool {
	nsec := true
	if withNsec != nil && withNsec[0] {
		nsec = t1.Nanosecond() == t2.Nanosecond()
	}
	return HourAndMinuteIsEqual(t1, t2) && SecondIsEqual(t1, t2) && nsec
}

//TimeIsEqualTo checks for hour, min, sec equality
//if with nsec is passed, will also compare nano secs
func TimeIsEqualTo(t1 time.Time, hour, min, sec int, nsec ...int) bool {
	nsecs := true
	if nsec != nil {
		nsecs = t1.Nanosecond() == nsec[0]
	}
	return t1.Hour() == hour && t1.Minute() == min && t1.Second() == sec && nsecs
}

//TimeIsGreaterThan checks if hour,min, sec are gt than time
func TimeIsGreaterThan(t1 time.Time, hour, min, sec int, nsec ...int) bool {
	nsecs := 0
	if nsec != nil {
		nsecs = nsec[0]
	}
	t2 := ReplaceTimeInDatetime(t1, hour, min, sec, nsecs)
	return t2.After(t1)
}

//TimeIsLessThan checks if hour,min, sec are gt than time
func TimeIsLessThan(t1 time.Time, hour, min, sec int, nsec ...int) bool {
	nsecs := 0
	if nsec != nil {
		nsecs = nsec[0]
	}
	t2 := ReplaceTimeInDatetime(t1, hour, min, sec, nsecs)
	return t2.Before(t1)
}

//DateIsEqualTo checks for hour, min, sec equality
//if with nsec is passed, will also compare nano secs
func DateIsEqualTo(t1 time.Time, year int, month time.Month, day int) bool {
	return t1.Day() == day && t1.Month() == month && t1.Year() == year
}

//HourAndMinuteIsEqual checks for hour min equality
func HourAndMinuteIsEqual(t1, t2 time.Time) bool {
	return (t1.Hour() == t2.Hour()) && (t1.Minute() == t2.Minute())
}

//DayIsEqual returns true if day field is same
func DayIsEqual(t1, t2 time.Time) bool {
	return t1.Day() == t2.Day()
}

//HourIsEqual returns true if hour field is same
func HourIsEqual(t1, t2 time.Time) bool {
	return t1.Hour() == t2.Hour()
}

//MinuteIsEqual returns true if minute field is same
func MinuteIsEqual(t1, t2 time.Time) bool {
	return t1.Minute() == t2.Minute()
}

//SecondIsEqual returns true if second field is same
func SecondIsEqual(t1, t2 time.Time) bool {
	return t1.Second() == t2.Second()
}

//NowTime returns only time.Now() time field
func NowTime() time.Time {
	return ExtractTimeFromDatetime(time.Now())
}

//NowDate returns only time.Now() date field
func NowDate() time.Time {
	return ExtractDateFromDatetime(time.Now())
}

//NowNextMinStart gets time.Now() and sets to next min start
//ie increment min by 1 and remove second, nsec
func NowNextMinStart() time.Time {
	return NowWithCustomTime(time.Now().Hour(), time.Now().Minute()+1, 0, 0)
}

//NowWithCustomTime uses todays date with time provided
func NowWithCustomTime(h, m, s int, nsec ...int) time.Time {
	nano := 0
	if nsec != nil {
		nano = nsec[0]
	}
	return ReplaceTimeInDatetime(NowDate(), h, m, s, nano)
}

//TimeIsInRange reports if t is between t1 and t2
//t1 is included and t2 is excluded
func TimeIsInRange(t, t1, t2 time.Time) bool {
	return !t.Before(t1) && t.Before(t2)
}

//TimeIsInArray checks if this exact time is in array
func TimeIsInArray(t time.Time, ta []time.Time) bool {
	for i := range ta {
		if t.Equal(ta[i]) {
			return true
		}
	}
	return false
}

//DurationDay returns time.Hour*24
func DurationDay() time.Duration {
	return time.Hour * 24
}

//DurationWeek returns time.Hour*24*7
func DurationWeek() time.Duration {
	return time.Hour * 24 * 7
}
