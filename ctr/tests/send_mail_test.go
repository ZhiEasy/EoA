package test

import (
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"testing"
)

func TestSendMail(t *testing.T) {
	//beego.AppConfig.String("ctr::email_username")
	//conf := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"host\":\"%s\",\"port\":%v}", beego.AppConfig.String("ctr::email_username"), beego.AppConfig.String("ctr::email_password"), beego.AppConfig.String("ctr::email_host"), beego.AppConfig.String("str::email_port"))
	conf := `{"username":"ahojcn@126.com","password":"WFENGENGEAEAQXAP","host":"smtp.126.com","port":25}`

	logrus.Errorln(conf)
	e := utils.NewEMail(conf)
	e.To = []string{"ahojcn@qq.com"}
	e.From = beego.AppConfig.String("ctr::email_username")
	e.Subject = "beego-邮件测试"
	e.Text = "text"
	e.HTML = "<h1>HTML!!!!</h1>"
	err := e.Send()
	logrus.Errorln(err)
}
