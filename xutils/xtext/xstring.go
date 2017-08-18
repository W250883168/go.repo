package xtext

import (
	"errors"
	"log"
	"runtime"
	"strings"

	"xutils/xapp"
)

// 测试是否为空字符串
func IsEmpty(str string) bool {
	return len(str) <= 0
}

// 测试是否为非空字符串
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// 测试是否为空白字符串
func IsBlank(str string) bool {
	val := strings.TrimSpace(str)
	return len(val) <= 0
}

// 测试是否为非空白字符串
func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

// 获取子字符串
func SubString(str string, index, length int) string {
	wchars := []rune(str)
	return string(wchars[index : index+length])
}

// 测试str为非空白字任串，否则引发panic
func RequireNonBlank(str string) {
	if IsBlank(str) {
		panic(errors.New("require non blank"))
	}
}

// 测试str1==str2, 否则引发panic
func RequireEqual(str1, str2 string) {
	if str1 != str2 {
		panic(errors.New("required str1 equals str2"))
	}
}

//消除空格
func SpaceTrim(str string) (str1 string) {
	str1 = strings.Replace(str, "\n", "", -1)
	str1 = strings.Replace(str1, " ", "", -1)
	return str1
}

// 例： kv:"PWR=00", split:"=", v:"00"； 返回："00"
func SplitValue(kv, split string) (v string) {
	if index := strings.Index(kv, split); index != -1 {
		wchars := []rune(kv)
		v = string(wchars[index+1:])
	}

	return v
}

// match="PWR={val}:",data="PWR=00:",split="{val}", 提取值为00
func FetchValue(match, data, split string) (v string) {
	cutstr := strings.Replace(match, split, "", 1)
	v = strings.Trim(data, cutstr)
	return v
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok && xapp.IsDebugMode() {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
