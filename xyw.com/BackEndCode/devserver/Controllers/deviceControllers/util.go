package deviceControllers

import (
	"gopkg.in/gorp.v1"

	userdao "dev.project/BackEndCode/devserver/DataAccess/usersDataAccess"
	"dev.project/BackEndCode/devserver/commons"
)

var gResponseMsgs = commons.ResponseMsgSet_Instance()

func doAuthValidate(tag string, roleType, userID int, dbmap *gorp.DbMap) bool {
	rd := userdao.CheckVaild(roleType, userID, tag, dbmap)
	valid := (rd.Rcode == "1000")
	return valid
}
