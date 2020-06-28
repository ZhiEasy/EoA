package main

import (
	_ "github.com/ahojcn/EoA/svr/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func main() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowAllOrigins:  true,
		AllowOrigins:     []string{"http://10.*.*.*:*", "http://localhost:*", "http://127.0.0.1:*", "http://*.*.*.*:*", "http://39.101.176.8:*", "*"},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
	}))

	// 返回主机基本信息
	//go func() {
	//	fileId, err := os.Open("id")
	//	if err != nil {
	//		logrus.Panicln("获取id失败", err.Error())
	//	}
	//	id, _, _ := bufio.NewReader(fileId).ReadLine()
	//	fileCb, err := os.Open("cb")
	//	if err != nil {
	//		logrus.Panicln("获取ctr位置失败", err.Error())
	//	}
	//	cb, _, _ := bufio.NewReader(fileCb).ReadLine()
	//	path := fmt.Sprintf("%s/cb/host/baseinfo", string(cb))
	//	logrus.Infof("id:%s, cb:%s, path:%s\n", id, cb, path)
	//
	//	type req struct {
	//		Id string `json:"id"`
	//		BaseInfo string `json:"base_info"`
	//	}
	//	d := req{
	//		Id:       string(id),
	//		BaseInfo: controllers.BaseInfo(),
	//	}
	//	bs, _ := json.Marshal(d)
	//	resp, err := http.Post(path, "application/json", bytes.NewReader(bs))
	//	if err != nil {
	//		logrus.Panicln("回调基本信息失败:", err.Error())
	//		return
	//	}
	//	logrus.Info(resp.Status)
	//}()

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
