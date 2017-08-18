package xnumeric

import (
	"errors"
	"fmt"
)

// 测试n, 必须满足(begin <= n < end), 否则引发panic
func RequireBetweenInt(n, begin, end int) {
	if (n < begin) || (n >= end) {
		panic(errors.New(fmt.Sprintf("Required Number In Range [%d,  %d), n=%d", begin, end, n)))
	}
}

// 测试n必须为非负数, 否则引发panic
func RequireNonNegative(n int) {
	if n < 0 {
		panic(errors.New(fmt.Sprintf("Required NonNegative Number, n=%d ", n)))
	}
}
