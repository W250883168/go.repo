package xobject

import (
	"errors"
	"fmt"

	"go.repo/xutils/xerr"
)

// 测试obj必须为非nil, 否则引发panic
func RequireNonNull(obj interface{}) {
	if obj == nil {
		xerr.ThrowPanic(errors.New("obj == nil"))
	}
}

// 获得字符串
func ToString(p interface{}) (str string) {
	if p != nil {
		str = stringize(p)
	}

	return str
}

// 获取格式化后字符串
func Fmt(p interface{}) string {
	return stringize(p)
}

func stringize(p interface{}) string {
	str := fmt.Sprint(p)
	if p == nil {
		return str
	}

	if obj, ok := p.(fmt.Stringer); ok {
		str = obj.String()
	} else if err, ok := p.(error); ok {
		str = err.Error()
	}

	return str
}
