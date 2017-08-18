package xhash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"runtime"

	"xutils/xapp"
)

type MD5 struct {
	reserved int
}

func NewMD5() *MD5 {
	return new(MD5)
}

func ToMD5(input string) (hash string) {
	h := md5.New()
	h.Write([]byte(input))
	hash = hex.EncodeToString(h.Sum(nil))
	return hash
}

func (p *MD5) MD5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok && xapp.IsDebugMode() {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
