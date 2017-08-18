package coapserver

import (
	"strings"
)

//根据匹配串匹配出变量值
//例子：paraMatchString="PWR={val}:",data="PWR=00:",varstr="{val}",匹配出来的值为00
func _MatchVarValue(paraMatchString, data string, varstr string) string {
	//log.Println("paraMatchString=", paraMatchString, " data=", data, " varstr=", varstr)
	if strings.Index(paraMatchString, varstr) == -1 { //如果不存在varstr，则返回空
		return ""
	}
	val := strings.Split(paraMatchString, varstr)
	rst := strings.Replace(data, val[0], "", 1)
	rst = strings.Replace(rst, val[1], "", 1)
	rst = strings.TrimSpace(rst)
	return rst
}
