package models

const (
	SUCCESS = 0
	NEED_UPDATE_INFO = 3002
	UNKNOW_ERROR = 4001
	REQUEST_ERROR = 4002
	AUTH_ERROR = 4003
	THIRED_ERROR = 4004
	HOST_CONN_ERROR = 4005
	PWD_ERROR = 4006
	HOST_REWATCH = 4007
	CPWD_ERROR = 4008
	SERVER_ERROR = 5001
)

var responseText = map[int]string {
	SUCCESS: "成功",
	NEED_UPDATE_INFO: "需要完善信息",
	UNKNOW_ERROR: "未知错误",
	REQUEST_ERROR: "请求错误",
	AUTH_ERROR: "未登录",
	SERVER_ERROR: "服务器错误",
	HOST_CONN_ERROR: "主机连接失败",
	PWD_ERROR: "密码错误",
	HOST_REWATCH: "重复关注",
	CPWD_ERROR: "密码不一致",
	THIRED_ERROR: "第三方系统错误",
}

func ResponseText(code int) string {
	str, ok := responseText[code]
	if ok {
		return str
	}
	return ResponseText(UNKNOW_ERROR)
}

type Response struct {
	Status int `json:"status"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}
