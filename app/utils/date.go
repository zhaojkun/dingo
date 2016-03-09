package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var conversion = map[rune]string{
	/*stdLongMonth      */ 'B': "January",
	/*stdMonth          */ 'b': "Jan",
	// stdNumMonth       */ 'm': "1",
	/*stdZeroMonth      */ 'm': "01",
	/*stdLongWeekDay    */ 'A': "Monday",
	/*stdWeekDay        */ 'a': "Mon",
	// stdDay            */ 'd': "2",
	// stdUnderDay       */ 'd': "_2",
	/*stdZeroDay        */ 'd': "02",
	/*stdHour           */ 'H': "15",
	// stdHour12         */ 'I': "3",
	/*stdZeroHour12     */ 'I': "03",
	// stdMinute         */ 'M': "4",
	/*stdZeroMinute     */ 'M': "04",
	// stdSecond         */ 'S': "5",
	/*stdZeroSecond     */ 'S': "05",
	/*stdLongYear       */ 'Y': "2006",
	/*stdYear           */ 'y': "06",
	/*stdPM             */ 'p': "PM",
	// stdpm             */ 'p': "pm",
	/*stdTZ             */ 'Z': "MST",
	// stdISO8601TZ      */ 'z': "Z0700",  // prints Z for UTC
	// stdISO8601ColonTZ */ 'z': "Z07:00", // prints Z for UTC
	/*stdNumTZ          */ 'z': "-0700", // always numeric
	// stdNumShortTZ     */ 'b': "-07",    // always numeric
	// stdNumColonTZ     */ 'b': "-07:00", // always numeric
	/* nonStdMilli		 */ 'L': ".000",
}

func DateFormat(t *time.Time, format string) string {
	retval := make([]byte, 0, len(format))
	for i, ni := 0, 0; i < len(format); i = ni + 2 {
		ni = strings.IndexByte(format[i:], '%')
		if ni < 0 {
			ni = len(format)
		} else {
			ni += i
		}
		retval = append(retval, []byte(format[i:ni])...)
		if ni+1 < len(format) {
			c := format[ni+1]
			if c == '%' {
				retval = append(retval, '%')
			} else {
				if layoutCmd, ok := conversion[rune(c)]; ok {
					retval = append(retval, []byte(t.Format(layoutCmd))...)
				} else {
					retval = append(retval, '%', c)
				}
			}
		} else {
			if ni < len(format) {
				retval = append(retval, '%')
			}
		}
	}
	return string(retval)
}

// Format unix time int64 to string
func DateInt64(ti int64, format string) string {
	t := time.Unix(int64(ti), 0)
	return DateTime(t, format)
}

// Format unix time string to string
func DateString(ts string, format string) string {
	i, _ := strconv.ParseInt(ts, 10, 64)
	return DateInt64(i, format)
}

// Format time.Time struct to string
// MM - month - 01
// M - month - 1, single bit
// DD - day - 02
// D - day 2
// YYYY - year - 2006
// YY - year - 06
// HH - 24 hours - 03
// H - 24 hours - 3
// hh - 12 hours - 03
// h - 12 hours - 3
// mm - minute - 04
// m - minute - 4
// ss - second - 05
// s - second = 5
func DateTime(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}

// Get unix stamp int64 of now
func NowUnix() int64 {
	return time.Now().Unix()
}

func Now() *time.Time {
	t := time.Now()
	return &t
}
