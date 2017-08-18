package vodControl

import (
	"strconv"

	"gopkg.in/gorp.v1"

	"dev.project/BackEndCode/devserver/DataAccess/usersDataAccess"
	"dev.project/BackEndCode/devserver/commons"
	"dev.project/BackEndCode/devserver/model/core"
)

func doVerifyAuth(roleType, userID int, tag string, dbmap *gorp.DbMap) (ok bool) {
	retData := usersDataAccess.CheckVaild(roleType, userID, tag, dbmap)

	responses := commons.ResponseMsgSet_Instance()
	ok = retData.Rcode == strconv.Itoa(responses.SUCCESS.Code)
	return ok
}

func doVerifyAuth2(tokens core.BasicsToken, tag string, dbmap *gorp.DbMap) (ok bool) {
	roleID := tokens.Rolestype
	userID := tokens.Usersid

	retData := usersDataAccess.CheckVaild(roleID, userID, tag, dbmap)
	responses := commons.ResponseMsgSet_Instance()
	ok = retData.Rcode == strconv.Itoa(responses.SUCCESS.Code)
	return ok
}
