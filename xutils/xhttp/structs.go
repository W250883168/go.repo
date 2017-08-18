package xhttp

// 分页信息
type PageInfo struct {
	PageIndex int // 当前页(1-index)
	PageSize  int // 每页大小(1-index)
	RowTotal  int // 总数量(行)
}

// 认证信息
type AuthInfo struct {
	UserID      string // 用户ID
	Password    string // 通行码
	Certificate string // 证书
}

// 请求
type HttpRequest struct {
	Auth    AuthInfo
	Page    PageInfo
	Content interface{}
}

// 响应
type HttpResponse struct {
	Page    PageInfo // 分页信息
	Content interface{}
}
