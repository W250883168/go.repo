package usersControllers

import (
	//	"fmt"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"dev.project/BackEndCode/devserver/DataAccess/usersDataAccess"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/viewmodel"

	"github.com/gin-gonic/gin"
)

//获取我的消息
func GetMyMessageList(c *gin.Context) {
	var rd core.Returndata
	var lg core.BasicsToken
	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "usersControllers|GetMyMessageList|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &pg)
	core.CheckErr(errs1, "usersControllers|GetMyMessageList|分页数据json格式转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.QueryMyMessage(lg, pg, dbmap)
		} else {
			rd.Rcode = "1004"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1003"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取我的消息
func GetIsNewMessageList(c *gin.Context) {
	var rd core.Returndata
	var lg core.BasicsToken
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "usersControllers|GetIsNewMessageList|登录数据json格式转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.QueryIsNewMessage(lg, dbmap)
		} else {
			rd.Rcode = "1004"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1003"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取学生信息
func GetStudentsinfo(c *gin.Context) {
	var rd core.Returndata
	var bt core.BasicsToken
	var lg viewmodel.GetStudentsinfo
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "usersControllers|GetStudentsinfo|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "usersControllers|GetStudentsinfo|令牌数据格式转换失败")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			ssarr := usersDataAccess.GetStudentsinfo(lg, bt, dbmap)
			if ssarr.Studentsid > 0 {
				rd.Result = ssarr
				rd.Rcode = "1000"
				rd.Reason = "成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "没有查询相关数据"
			}
		} else {
			rd.Rcode = "1002"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}
