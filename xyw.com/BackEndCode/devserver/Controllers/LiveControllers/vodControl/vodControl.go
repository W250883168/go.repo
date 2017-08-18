package vodControl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"

	liveProvider "dev.project/BackEndCode/devserver/DataAccess/liveDataAccess"
	"dev.project/BackEndCode/devserver/commons"
	"dev.project/BackEndCode/devserver/commons/xdebug"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/viewmodel/videoview"
)

// 查询视频列表
func GetVideoList(c *gin.Context) {
	var pd core.PageData
	var rd = core.Returndata{Result: &pd}
	responses := commons.ResponseMsgSet_Instance()
	var tokens core.BasicsToken
	var err error
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err = json.Unmarshal(data, &tokens); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	TAG := "getquerylivelist"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doVerifyAuth2(tokens, TAG, dbmap); !ok {
		rd.Rcode = strconv.Itoa(responses.ROLE_AUTH_LIMITED.Code)
		rd.Reason = responses.ROLE_AUTH_LIMITED.Text
		c.JSON(http.StatusOK, rd)
		return
	}

	var request videoview.Request_VideoInfoView
	if err = json.Unmarshal(data, &request); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	rd.Rcode = strconv.Itoa(responses.FOUND_NODATA.Code)
	rd.Reason = responses.FOUND_NODATA.Text
	if list, ok := liveProvider.QueryVideoInfos(&request, dbmap); ok {
		rd.Rcode = strconv.Itoa(responses.SUCCESS.Code)
		rd.Reason = responses.SUCCESS.Text
		pd.PageData = list
		pd.PageSize = request.PageInfo.PageSize
		pd.PageIndex = request.PageInfo.PageIndex
		pd.PageCount = request.PageInfo.RowTotal
	}

	c.JSON(http.StatusOK, rd)
}

// 获取视频详细
func GetVideoDetails(c *gin.Context) {
	var rd = core.Returndata{Result: videoview.VideoDetailView{}}
	responses := commons.ResponseMsgSet_Instance()
	var tokens core.BasicsToken
	var err error
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err = json.Unmarshal(data, &tokens); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	var request videoview.VideoDetailView
	if err = json.Unmarshal(data, &request); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	rd.Rcode = strconv.Itoa(responses.FOUND_NODATA.Code)
	rd.Reason = responses.FOUND_NODATA.Text
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if video, ok := liveProvider.QueryVideoDetials(request.ID, dbmap); ok {
		rd.Rcode = strconv.Itoa(responses.SUCCESS.Code)
		rd.Reason = responses.SUCCESS.Text
		rd.Result = video
	}

	c.JSON(http.StatusOK, rd)
}

// 视频详细修改
func UpdateVideoDetatils(c *gin.Context) {
	var rd = core.Returndata{Result: ""}
	responses := commons.ResponseMsgSet_Instance()
	var tokens core.BasicsToken
	var err error
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err = json.Unmarshal(data, &tokens); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	var request videoview.VideoDetailView
	if err = json.Unmarshal(data, &request); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	rd.Rcode = strconv.Itoa(responses.FAIL.Code)
	rd.Reason = responses.FAIL.Text
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := liveProvider.UpdateVideoDetails(request, dbmap); ok {
		rd.Rcode = strconv.Itoa(responses.SUCCESS.Code)
		rd.Reason = responses.SUCCESS.Text
	}

	c.JSON(http.StatusOK, rd)
}

// 视频删除
func DeleteVideo(c *gin.Context) {
	var rd = core.Returndata{Result: ""}
	responses := commons.ResponseMsgSet_Instance()
	var tokens core.BasicsToken
	var err error
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err = json.Unmarshal(data, &tokens); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	var request videoview.VideoDetailView
	if err = json.Unmarshal(data, &request); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	rd.Rcode = strconv.Itoa(responses.FAIL.Code)
	rd.Reason = responses.FAIL.Text
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := liveProvider.DeleteVideo(request.ID, dbmap); ok {
		rd.Rcode = strconv.Itoa(responses.SUCCESS.Code)
		rd.Reason = responses.SUCCESS.Text
	}

	c.JSON(http.StatusOK, rd)
}

// 删除视频附件
func DeleteAttachment(c *gin.Context) {
	var rd = core.Returndata{Result: ""}
	responses := commons.ResponseMsgSet_Instance()
	var tokens core.BasicsToken
	var err error
	data, err := ioutil.ReadAll(c.Request.Body)
	if err = json.Unmarshal(data, &tokens); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	var request videoview.AttachmentView
	if err = json.Unmarshal(data, &request); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		xdebug.DebugError(err)
		c.JSON(http.StatusOK, rd)
		return
	}

	rd.Rcode = strconv.Itoa(responses.FAIL.Code)
	rd.Reason = responses.FAIL.Text
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := liveProvider.DeleteAttachment(request.ID, dbmap); ok {
		rd.Rcode = strconv.Itoa(responses.SUCCESS.Code)
		rd.Reason = responses.SUCCESS.Text
	}

	c.JSON(http.StatusOK, rd)
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	// data, _ := json.Marshal(videoview.VideoDetailView{})
	// fmt.Printf("%+v\n", string(data))
}
