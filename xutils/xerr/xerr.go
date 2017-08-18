package xerr

import (
	"fmt"
	"xutils/xdebug"
)

// 捕获异常
func Catch() {
	if catch := recover(); catch != nil {
		switch catch.(type) {
		case error:
			xdebug.LogError(catch.(error))
		case fmt.Stringer:
			str := catch.(fmt.Stringer)
			xdebug.LogString(str.String())
		default:
			xdebug.LogString(fmt.Sprint(catch))
		}
	}
}

// 捕获异常(error)
func CatchPanic() {
	if catch := recover(); catch != nil {
		err := catch.(error)
		xdebug.LogError(err)
	}
}

// 抛出异常(error)
func ThrowPanic(err error) {
	if err != nil {
		panic(err)
	}
}
