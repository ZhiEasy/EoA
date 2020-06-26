package models

const (
	SUCCESS            = 0
	UNKNOW_ERROR       = 4001
)

var responseText = map[int]string{
	SUCCESS:            "成功",
	UNKNOW_ERROR:       "未知错误",
}

func ResponseText(code int) string {
	str, ok := responseText[code]
	if ok {
		return str
	}
	return ResponseText(UNKNOW_ERROR)
}

type Response struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}
