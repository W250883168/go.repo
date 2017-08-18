package xtext

import "strings"

const (
	_BLANK_STRING string = " \t\r\n"
)

func IsEmpty(str string) bool {
	return len(str) <= 0
}

func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

func IsBlank(str string) bool {
	val := strings.Trim(str, _BLANK_STRING)
	return len(val) <= 0
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}
