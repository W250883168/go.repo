package enclosure

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"basicproject/commons"
	"basicproject/commons/xdebug"
	"basicproject/commons/xfile"
	"basicproject/model/curriculum"
	"basicproject/viewmodel/videoview"
	core "xutils/xcore"
)

// 接受上传附件
func UploadAttachment(c *gin.Context) {
	defer xdebug.DoRecover() // 错误恢复
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	responses := commons.ResponseMsgSet_Instance()
	var rd = core.Returndata{
		Rcode:  strconv.Itoa(responses.FAIL.Code),
		Reason: responses.FAIL.Text,
		Result: ""}
	defer func() { c.JSON(http.StatusOK, rd) }()

	_, _, err := c.Request.FormFile("file")
	xdebug.HandleError(err) //
	if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
		for _, pFileHeaders := range c.Request.MultipartForm.File {
			if len(pFileHeaders) > 0 {
				var pFileHead *multipart.FileHeader = pFileHeaders[0]
				// fmt.Printf("%+v\n%+v\n", k, pFileHead)
				inFile, err := pFileHead.Open()
				defer inFile.Close()
				xdebug.HandleError(err)

				timestr := time.Now().Format("20060102")
				fpath := "/templates/upfile"
				dir := filepath.Join(xfile.GetPWD(), filepath.FromSlash(fpath), timestr)
				if !xfile.FileExist(dir) { // 判断文件目录是否存在
					err = os.Mkdir(dir, os.ModePerm) // 创建文件目录
					xdebug.HandleError(err)
				}

				fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + filepath.Ext(pFileHead.Filename)
				newFile, err := os.Create(dir + "/" + fileName)
				defer newFile.Close()
				xdebug.HandleError(err)

				rd.Rcode = strconv.Itoa(responses.DATA_VERIFY_FAIL.Code)
				rd.Reason = responses.DATA_VERIFY_FAIL.Text
				_, err = io.Copy(newFile, inFile)
				xdebug.HandleError(err)

				httpVirtualPath := path.Join("/web/upfile", timestr, fileName)
				idValue := c.Request.FormValue("CurriculumClassroomChapterID")
				id, err := strconv.Atoi(idValue)
				var attach = curriculum.Enclosure{
					Enclosurename:        pFileHead.Filename,
					Enclosurepath:        filepath.Join(dir, fileName),
					EnclosureVirtualPath: httpVirtualPath,
					Createdate:           time.Now().Format("2016-08-27 15:21:03"),
					Enclosuretype:        filepath.Ext(pFileHead.Filename),
					// Enclosuresize
					// Enclosureicon
					// IsPublish
					Curriculumclassroomchaptercentreid: id}
				dbmap := core.InitDb()
				defer dbmap.Db.Close()
				dbmap.AddTableWithName(curriculum.Enclosure{}, "enclosure").SetKeys(true, "Id")
				err = dbmap.Insert(&attach)
				xdebug.HandleError(err)

				// fmt.Printf("%+v\n", attach)
				rd.Rcode = strconv.Itoa(responses.SUCCESS.Code)
				rd.Reason = responses.SUCCESS.Text
				rd.Result = videoview.AttachmentView{
					ID:                           attach.Id,
					EnclosureName:                attach.Enclosurename,
					EnclosureType:                attach.Enclosuretype,
					EnclosureSize:                attach.Enclosuresize,
					VirtualPath:                  attach.EnclosureVirtualPath,
					CreateDate:                   attach.Createdate,
					IsPublish:                    attach.IsPublish,
					EnclosuerIcon:                attach.Enclosureicon,
					CurriculumClassroomChapterID: attach.Curriculumclassroomchaptercentreid}
				break
			}
		}
	}
}

//func UploadAttachment2(c *gin.Context) {
//	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
//	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

//	var err error
//	responses := commons.ResponseMsgSet_Instance()
//	var rd core.Returndata
//	fmt.Printf("	<<<<<<<<<<Request: \n%+v\n", c.Request)
//	fmt.Printf("	<<<<<<<<<<Request.Header: \n%+v\n", c.Request.Header)
//	fmt.Printf("	<<<<<<<<<<Context.Keys: \n%+v\n", c.Keys)
//	fmt.Printf("	<<<<<<<<<<Context.Params: \n%+v\n", c.Params)
//	fmt.Printf("	<<<<<<<<<<Request.Form: \n%+v\n", c.Request.Form["upfile"])
//	fmt.Printf("	<<<<<<<<<<Request.PostForm: \n%+v\n", c.Request.PostForm["upfile"])

//	timestr := time.Now().Format("20060102")
//	fpath := "/templates/upfile"
//	dir := filepath.Join(xfile.GetPWD(), filepath.FromSlash(fpath), timestr)
//	if !xfile.FileExist(dir) { // 判断文件目录是否存在
//		err = os.Mkdir(dir, os.ModePerm) // 创建文件目录
//		xdebug.HandleError(err)
//	}

//	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ".png"
//	newFile, err := os.Create(dir + "/" + fileName)
//	defer newFile.Close()
//	xdebug.DebugError(err)

//	rd.Rcode = strconv.Itoa(responses.DATA_VERIFY_FAIL.Code)
//	rd.Reason = responses.DATA_VERIFY_FAIL.Text
//	_, err = io.Copy(newFile, c.Request.Body)
//	xdebug.DebugError(err)
//	if err == nil {
//		rd.Rcode = strconv.Itoa(responses.SUCCESS.Code)
//		rd.Reason = responses.SUCCESS.Text
//		rd.Result = "/web/upfile/" + timestr + "/" + fileName

//	}

//	c.JSON(http.StatusOK, rd)
//}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

}
