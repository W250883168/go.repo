package httpview

import (
	"fmt"
	"log"
	"runtime"
)

// 分页信息
type PageInfoView struct {
	PageIndex int
	PageSize  int
	RowTotal  int
}

// HTTP请求
type HttpRequestView struct {
	PageInfoView
	Data interface{}
}

// HTTP响应
type HttpResponseView struct {
	PageInfoView
	HttpCode int
	Content  interface{}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	fmt.Print()
}
