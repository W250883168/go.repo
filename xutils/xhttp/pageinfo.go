package xhttp

import (
	"errors"
	"fmt"
)

// 获得总页数
func (p *PageInfo) PageTotal() (total int) {
	if p.Valid() {
		total = (p.RowTotal / p.PageSize)
		if p.RowTotal%p.PageSize != 0 {
			total++
		}
	}

	return total
}

// SQL分页查询条件
func (p *PageInfo) SQL_LimitString() (limit string) {
	if p.Valid() {
		offset := (p.PageIndex - 1) * p.PageSize
		count := p.PageSize
		limit = fmt.Sprintf(`LIMIT %d, %d `, offset, count)
	}

	return limit
}

// 检查数据有效性, 无效引发Panic
func (p *PageInfo) Validate() {
	if !p.Valid() {
		panic(errors.New("(PageIndex < 1) OR (PageSize < 1)"))
	}
}

// 检查数据有效性
func (p *PageInfo) Valid() (ok bool) {
	ok = (p.PageIndex >= 1) && (p.PageSize >= 1)
	return ok
}
