package xtime

import (
	"fmt"
	"time"
)

type TimeSpan time.Duration

func (p *TimeSpan) String() string {
	dur := time.Duration(*p)

	hh := time.Duration(dur.Hours())
	dur -= time.Hour * hh
	mm := time.Duration(dur.Minutes())
	dur -= time.Minute * mm
	ss := time.Duration(dur.Seconds())
	dur -= time.Second * ss
	ms := dur / time.Millisecond

	var args = []interface{}{hh, mm, ss}
	format := "%0.2d:%0.2d:%0.2d"
	if ms != 0 {
		args = append(args, ms)
		format += ".%d"
	}

	return fmt.Sprintf(format, args...)
}

// 获得指定格式的字符串文本(如:format=[%0.2d小时%0.2d分钟%0.2d秒])
func (p *TimeSpan) ToString(format string) string {
	dur := time.Duration(*p)

	hh := time.Duration(dur.Hours())
	dur -= time.Hour * hh
	mm := time.Duration(dur.Minutes())
	dur -= time.Minute * mm
	ss := time.Duration(dur.Seconds())
	dur -= time.Second * ss

	var args = []interface{}{hh, mm, ss}
	return fmt.Sprintf(format, args...)
}
