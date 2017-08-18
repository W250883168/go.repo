package main_action

import (
	"TimingService/Controllers"
	"TimingService/Service"

	"github.com/gin-gonic/gin"
)

func LoadWebFile(r *gin.Engine) {
}
func LoadAction(b *gin.Engine) {
	r := b.Group("/Task", func(c *gin.Context) {})
	r.POST("/AddTimedTask", timingControllers.SaveTimedTask)             //添加定时任务
	r.POST("/ChangeTimedTask", timingControllers.ChangeTimedTask)        //修改定时任务
	r.POST("/DelTimedTask", timingControllers.DelTimedTask)              //删除定时任务
	r.POST("/TimedTaskList", timingControllers.GetTimedTasklist)         //查看定时任务
	r.POST("/TimedTaskInfo", timingControllers.GetTimedTaskinfo)         //查看定时任务详细情况
	r.POST("/OnOrOffTimedTask", timingControllers.PostOnOrOffTimedTask)  //关闭/打开定时任务
	r.POST("/EventSetTableList", timingControllers.GetEventSetTablelist) //查看响应事件列表
	go timingService.RunExecCmd()
	go timingService.RunUpdateDate()
	go timingService.RunBackgroundServer()
	go timingService.RunBackgroundCheckValid()
}
