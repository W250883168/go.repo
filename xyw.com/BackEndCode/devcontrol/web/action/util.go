package action

import (
	"strings"
	"time"
)

//判断是否在线
//from:起始时间
//duration:时间段，如30s,1d等
//如果from + duration的时间还小于当前时间，则视为离线
func _IsOnline(from, duration string) (result bool) {
	from1, err := time.Parse("2006-01-02 15:04:05", from)
	d, _ := time.ParseDuration(duration)
	from2 := from1.Add(d) //from + duration 后的时间

	time2 := time.Now().Format("2006-01-02 15:04:05")
	now, err := time.Parse("2006-01-02 15:04:05", time2) //现在的时间

	if err == nil && from2.Before(now) {
		result = false
	} else {
		result = true
	}
	return
}

//根据匹配串匹配出变量值
//例子：paraMatchString="PWR={val}:",data="PWR=00:",varstr="{val}",匹配出来的值为00
func _MatchVarValue(paraMatchString, data, varstr string) string {
	if strings.Index(paraMatchString, varstr) == -1 { //如果不存在varstr，则返回空
		return ""
	}
	val := strings.Split(paraMatchString, varstr)
	rst := strings.Replace(data, val[0], "", 1)
	rst = strings.Replace(rst, val[1], "", 1)
	rst = strings.TrimSpace(rst)
	return rst
}

// 获得命令码显示文本
func _CmdCode_ToName(cmd, defaultName string) (str string) {
	str = defaultName
	if cmd == "on" {
		str = "开启"
	} else if cmd == "off" {
		str = "关闭"
	} else if cmd == "toggle" {
		str = "切换端口"
	}

	return str
}
