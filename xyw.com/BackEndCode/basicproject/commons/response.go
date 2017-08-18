package commons

import "strconv"

var pInstance_ResponseMsgSet *ResponseMsgSet

type ResponseMsg struct {
	Code int    // 响应码
	Text string // 说明文本
}

type ResponseMsgSet struct {
	SUCCESS           ResponseMsg
	FAIL              ResponseMsg
	DATA_VERIFY_FAIL  ResponseMsg
	DATA_MALFORMED    ResponseMsg
	TOKEN_INCORRECT   ResponseMsg
	AUTH_LIMITED      ResponseMsg
	ROLE_AUTH_LIMITED ResponseMsg
	EXE_CMD_FAIL      ResponseMsg
	SEND_UDP_FAIL     ResponseMsg
	FOUND_NODATA      ResponseMsg
}

func ResponseMsgSet_Instance() *ResponseMsgSet {
	return pInstance_ResponseMsgSet
}

func (p *ResponseMsg) CodeText() string {
	return strconv.Itoa(p.Code)
}

func init() {
	pInstance_ResponseMsgSet = &ResponseMsgSet{
		SUCCESS:           ResponseMsg{Code: 1000, Text: "操作成功!"},
		FAIL:              ResponseMsg{Code: 1001, Text: "操作失败!"},
		DATA_VERIFY_FAIL:  ResponseMsg{Code: 1002, Text: "数据验证失败!"},
		DATA_MALFORMED:    ResponseMsg{Code: 1003, Text: "数据格式无效!"},
		TOKEN_INCORRECT:   ResponseMsg{Code: 1004, Text: "令牌无效!"},
		AUTH_LIMITED:      ResponseMsg{Code: 1005, Text: "此功能未授权!"},
		ROLE_AUTH_LIMITED: ResponseMsg{Code: 1006, Text: "角色权限受限!"},
		EXE_CMD_FAIL:      ResponseMsg{Code: 1007, Text: "执行命令失败!"},
		SEND_UDP_FAIL:     ResponseMsg{Code: 1008, Text: "发送UDP请求失败!"},
		FOUND_NODATA:      ResponseMsg{Code: 1099, Text: "未找到数据!"}}

	// fmt.Println(pInstance_ResponseMsgSet)
}
