package controllers

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
)

// EmailLogController operations for EmailLog
type EmailLogController struct {
	BaseController
}

func SendMail(emailList []string, subject string, html string) {
	go func() {
		conf := fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\",\"host\":\"%s\",\"port\":%v}", beego.AppConfig.String("ctr::email_username"), beego.AppConfig.String("ctr::email_password"), beego.AppConfig.String("ctr::email_host"), beego.AppConfig.String("ctr::email_port"))
		logrus.Warnln(conf)
		e := utils.NewEMail(conf)
		e.To = emailList
		e.From = beego.AppConfig.String("ctr::email_username")
		e.Subject = subject
		e.HTML = html
		err := e.Send()
		if err != nil {
			logrus.Errorf("发送邮件失败\nemailList:%v\n错误：%v", emailList, err)
		}
	}()
}
