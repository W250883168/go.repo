package xtime

import (
	"time"

	"github.com/jinzhu/now"
)

const gFormatString = "2006-01-02 15:04:05"
const gFormatString2 = "20060102-150405.999"
const FORMAT_yyyyMMddHHmmss = "20060102-150405"
const FORMAT_yyyyMMddHHmmssfff = "20060102-150405.999"
const FORMAT_yyyyMM = "2006-01"

func TimeString(t *time.Time) (str string) {
	if t != nil {
		str = t.Format(gFormatString)
	}

	return str
}

func NowString() string {
	return time.Now().Format(gFormatString)
}

func FormatString() string {
	return gFormatString
}

func Parse(str string) (tm time.Time) {
	tm, _ = time.Parse(gFormatString, str)
	return tm
}

func MonthDays(t *time.Time) int {
	tEnd := now.New(*t).EndOfMonth()
	_, _, days := tEnd.Date()
	return days
}

func WeekdayName(p *time.Time) (str string) {
	switch p.Weekday() {
	case time.Monday:
		str = "星期一"
	case time.Tuesday:
		str = "星期二"
	case time.Wednesday:
		str = "星期三"
	case time.Thursday:
		str = "星期四"
	case time.Friday:
		str = "星期五"
	case time.Saturday:
		str = "星期六"
	case time.Sunday:
		str = "星期天"
	}

	return str
}

func GetChainDay(tm time.Time) (str string) {
	switch tm.Weekday().String() {
	case "Monday":
		str = "星期一"
	case "Tuesday":
		str = "星期二"
	case "Wednesday":
		str = "星期三"
	case "Thursday":
		str = "星期四"
	case "Friday":
		str = "星期五"
	case "Saturday":
		str = "星期六"
	case "Sunday":
		str = "星期天"
	}

	return str
}
