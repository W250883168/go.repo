package users

import (
	"bytes"
	"fmt"
	"dev.project/BackEndCode/devserver/commons/xtext"
)

type QueryWhere4User struct {
	Searhtxt  string // 关键词
	RoleID    int    // 角色类型
	Sex       int    // 性别
	State     int    // 状态
	PageIndex int    // 页面索引
	PageSize  int    // 页面大小
}

func (user *QueryWhere4User) WhereString() string {
	buff := bytes.NewBufferString(" WHERE (1=1) ")
	if user != nil {
		if xtext.IsNotBlank(user.Searhtxt) {
			str := percentstr(user.Searhtxt)
			txt := fmt.Sprintf(` AND ((TUser.Loginuser LIKE '%s') OR (TUser.Truename LIKE '%s') OR (TUser.Nickname LIKE '%s')) `, str, str, str)
			buff.WriteString(txt)
		}

		if user.RoleID > 0 {
			str := fmt.Sprintf(` AND (TUser.Rolesid = %d) `, user.RoleID)
			buff.WriteString(str)
		}
		if user.RoleID == -1 {
			str := ` AND (TUser.Rolesid not in(2,3)) `
			buff.WriteString(str)
		}

		if user.State > 0 {
			str := fmt.Sprintf(` AND (TUser.Userstate = %d) `, user.State)
			buff.WriteString(str)
		}
	}

	return buff.String()
}

func (user *QueryWhere4User) LimitString() string {
	offset := (user.PageIndex - 1) * user.PageSize
	count := user.PageSize
	limit_str := fmt.Sprintf(`LIMIT %d, %d `, offset, count)
	return limit_str

}

type QueryCondition_Teacher struct {
	QueryWhere4User
}

type QueryCondition_Student struct {
	QueryWhere4User
}

func percentstr(txt string) (ret string) {
	if len(txt) > 0 {
		ret = "%" + txt + "%"
	}

	return ret
}
