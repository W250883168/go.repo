package liveControllers

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"dev.project/BackEndCode/devserver/DataAccess/liveDataAccess"
	"dev.project/BackEndCode/devserver/DataAccess/usersDataAccess"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/viewmodel"

	"github.com/gin-gonic/gin"
)

/*
按条件查询课程直播或点播相关数据
*/
func GetQueryLiveList(c *gin.Context) {
	var rd core.Returndata
	var lgs viewmodel.QueryLivesWhere
	var lg core.BasicsToken
	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "liveControllers|GetQueryLiveList|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &lgs)
	core.CheckErr(errs1, "liveControllers|GetQueryLiveList|数据转换错误")
	errs1 = json.Unmarshal(data, &pg)
	core.CheckErr(errs1, "liveControllers|GetQueryLiveList|分页数据转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "getquerylivelist", dbmap)
			if rd.Rcode == "1000" {
				lgs.Psstate = -1
				lgs.Trstate = -1
				res := liveDataAccess.QueryLiveWhereAll(lgs, pg, dbmap)
				rd.Result = res
				rd.Rcode = "1000"
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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

/*
查询上周视频
*/
func GetQueryLastWeeklivelist(c *gin.Context) {
	var rd core.Returndata
	var lgs viewmodel.QueryLivesWhere
	var lg core.BasicsToken
	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "liveControllers|GetQueryLastWeeklivelist|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &lgs)
	core.CheckErr(errs1, "liveControllers|GetQueryLastWeeklivelist|数据转换错误")
	errs1 = json.Unmarshal(data, &pg)
	core.CheckErr(errs1, "liveControllers|GetQueryLastWeeklivelist|分页数据转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetQueryLastWeeklivelist", dbmap)
			if rd.Rcode == "1000" {
				lgs.Psstate = -1
				lgs.Trstate = -1
				res := liveDataAccess.QueryLiveWhereAll(lgs, pg, dbmap)
				for k := 0; k < len(res); k++ {
					lgs.Curriculumsid = res[k].Curriculumsid
					res[k].Lives = liveDataAccess.QueryLiveInfo(lgs, dbmap)
				}
				rd.Result = res
				rd.Rcode = "1000"
				rd.Reason = ""
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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

/*
查询我缺课的课程列表
*/
func GetMyAbsentList(c *gin.Context) {
	var rd core.Returndata
	var lgs viewmodel.QueryLivesWhere
	var lg core.BasicsToken
	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "liveControllers|GetMyAbsentList|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &lgs)
	core.CheckErr(errs1, "liveControllers|GetMyAbsentList|数据转换错误")
	errs1 = json.Unmarshal(data, &pg)
	core.CheckErr(errs1, "liveControllers|GetMyAbsentList|分页数据转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			lgs.Trstate = 2
			lgs.Psstate = 0
			lgs.PsUserid = lg.Usersid
			res := liveDataAccess.QueryMyAbsentList(lgs, pg, dbmap)
			for k := 0; k < len(res); k++ {
				lgs.Curriculumsid = res[k].Curriculumsid
				res[k].Lives = liveDataAccess.QueryMyAbsentInfo(lgs, dbmap)
			}
			rd.Result = res
			rd.Rcode = "1000"
			rd.Reason = ""
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

/*
查询我缺课的课程章节列表详细信息
*/
func GetMyAbsentInfo(c *gin.Context) {
	var rd core.Returndata
	var lgs viewmodel.QueryLivesWhere
	var lg core.BasicsToken
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "liveControllers|GetMyAbsentInfo|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &lgs)
	core.CheckErr(errs1, "liveControllers|GetMyAbsentInfo|数据转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			lgs.Trstate = 2
			lgs.Psstate = 0
			lgs.PsUserid = lg.Usersid
			res := liveDataAccess.QueryMyAbsentInfo(lgs, dbmap)
			rd.Result = res
			rd.Rcode = "1000"
			rd.Reason = ""
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

/*
查询课程下的所有章节播放列表
*/
func GetQueryLiveInfo(c *gin.Context) {
	var rd core.Returndata
	var lgs viewmodel.QueryLivesWhere
	var lg core.BasicsToken
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "liveControllers|GetQueryLiveInfo|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &lgs)
	core.CheckErr(errs1, "liveControllers|GetQueryLiveInfo|数据转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			lgs.Psstate = -1
			lgs.Trstate = -1
			res := liveDataAccess.QueryLiveInfo(lgs, dbmap)
			rd.Result = res
			rd.Rcode = "1000"
			rd.Reason = ""
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

/*
查询课程下的所有章节播放列表
*/
func GetQueryLiveEnclosure(c *gin.Context) {
	var rd core.Returndata
	var lgs viewmodel.QueryLivesWhere
	var lg core.BasicsToken
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "liveControllers|GetQueryLiveEnclosure|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &lgs)
	core.CheckErr(errs1, "liveControllers|GetQueryLiveEnclosure|数据转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			res := liveDataAccess.QueryLiveEnclosure(lgs, dbmap)
			rd.Result = res
			rd.Rcode = "1000"
			rd.Reason = ""
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

/*
更新播放次数
*/
func UpdateLivePlayNum(c *gin.Context) {
	var rd core.Returndata
	var lgs viewmodel.QueryLivesWhere
	var lg core.BasicsToken
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "liveControllers|UpdateLivePlayNum|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &lgs)
	core.CheckErr(errs1, "liveControllers|UpdateLivePlayNum|数据转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			lgs.PsUserid = lg.Usersid
			res := liveDataAccess.UpdateLivePlayNum(lgs, dbmap)
			if res.Rcode == "" {
				rd.Rcode = "1000"
				rd.Reason = ""
			} else {
				rd = res
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

/*
更新下载次数
*/
func UpdateLiveDownloadNum(c *gin.Context) {
	var rd core.Returndata
	var lgs viewmodel.QueryLivesWhere
	var lg core.BasicsToken
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "liveControllers|UpdateLiveDownloadNum|登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &lgs)
	core.CheckErr(errs1, "liveControllers|UpdateLiveDownloadNum|数据转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			lgs.PsUserid = lg.Usersid
			res := liveDataAccess.UpdateLiveDownloadNum(lgs, dbmap)
			if res.Rcode == "" {
				rd.Rcode = "1000"
				rd.Reason = ""
			} else {
				rd = res
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
