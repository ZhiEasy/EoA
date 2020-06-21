package models

const (
	SUCCESS = 0
	UNKNOW_ERROR = 4001
	REQUEST_ERROR = 4002
	AUTH_ERROR = 4003
	THIRED_ERROR = 4004
	SERVER_ERROR = 5001
)

var responseText = map[int]string {
	SUCCESS: "成功",
	UNKNOW_ERROR: "未知错误",
	REQUEST_ERROR: "请求错误",
	AUTH_ERROR: "权限错误",
	SERVER_ERROR: "服务器错误",
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
