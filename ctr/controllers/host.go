package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ahojcn/EoA/ctr/models"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

// HostController operations for Host
type HostController struct {
	BaseController
}

// 添加主机
func (c *HostController)AddHost() {
	userId := c.LoginRequired(true)
	userObj, err := models.GetUserById(userId)
	if err != nil {
		c.ReturnResponse(models.AUTH_ERROR, nil, true)
	}

	var req models.AddHostReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.ReturnResponse(models.REQUEST_ERROR, nil, true)
	}

	if req.Check() == false {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}

	// 测试主机是否可以连接
	cliConf := SSHClientConfig{
		Host:       req.Ip,
		Port:       22,
		Username:   req.LoginName,
		Password:   req.LoginPwd,
	}
	start := time.Now().UnixNano()
	err = cliConf.CreateClient()
	end := time.Now().UnixNano()
	if err != nil {
		c.ReturnResponse(models.HOST_CONN_ERROR, nil, true)
	}

	hash := md5.New()
	hash.Write([]byte(req.LoginPwd))
	hostObj := models.Host{
		UserId:        userObj,
		Ip:            req.Ip,
		Name:          req.Name,
		Description:   req.Description,
		LoginName:     req.LoginName,
		LoginPwd:      hex.EncodeToString(hash.Sum(nil)),
	}
	//hostId, err := models.AddHost(&hostObj)
	_, err = models.AddHost(&hostObj)
	if err != nil {
		logrus.Warningf("User:%v 添加主机失败，Request：%v，错误信息：%v", userId, req, err)
		c.ReturnResponse(models.SERVER_ERROR, nil, true)
	}

	// TODO 开一个协程去获取主机 base info
	// TODO 使用回调的方式获取基本信息
	//go func() {
	//	shell := fmt.Sprintf("mkdir -p /data/eoa/conf && cd /data/eoa/ && wget %s && cd conf && wget %s && cd .. && chmod +x svr_linux_amd64 && nohup ./svr_linux_amd64 > svr.log 2>&1 &", beego.AppConfig.String("svr::svrbinpath"), beego.AppConfig.String("svr::svrconfpath"))
	//	logrus.Errorln("执行脚本: ", shell)
	//	err := cliConf.RunShell(shell)
	//	logrus.Errorln("执行脚本结果: ", cliConf.LastResult)
	//	if err != nil {
	//		logrus.Errorln("执行脚本出错: ", err.Error())
	//		hostObj.BaseInfo = fmt.Sprintf("获取信息失败:%v", err.Error())
	//		_ = models.UpdateHostById(&hostObj)
	//	}
	//	httpCli := http.Client{}
	//	resp, err := httpCli.Get(fmt.Sprintf("http://%s:%s%s", req.Ip, beego.AppConfig.String("svr::svrport"), beego.AppConfig.String("svr::svrbaseinfopath")))
	//	if err != nil {
	//		logrus.Errorln("获取信息失败: ", err.Error())
	//		hostObj.BaseInfo += fmt.Sprintf("获取信息失败:%v", err.Error())
	//		_ = models.UpdateHostById(&hostObj)
	//		return
	//	}
	//	info, _ :=ioutil.ReadAll(resp.Body)
	//	logrus.Errorln(info)
	//}()

	// 重复关注，不用管
	_ = AddHostWatch(userObj.Id, hostObj.Id)
	//if err != nil {
	//	c.ReturnResponse(models.HOST_REWATCH, nil, true)
	//}

	d := make(map[string]string)
	d["used"] = fmt.Sprintf("连接用时 %v ms", (end - start)/1e6)
	c.ReturnResponse(models.SUCCESS, d, true)
}

// 测试主机连接
func (c *HostController) HostConnectionTest()  {
	c.LoginRequired(true)

	var req models.HostConnection
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.ReturnResponse(models.REQUEST_ERROR, nil, true)
	}

	cliConf := SSHClientConfig{
		Host:       req.Ip,
		Port:       22,
		Username:   req.LoginName,
		Password:   req.LoginPwd,
	}
	logrus.Warnln(req)
	start := time.Now().UnixNano()
	err = cliConf.CreateClient()

	end := time.Now().UnixNano()
	if err != nil {
		c.ReturnResponse(models.HOST_CONN_ERROR, nil, true)
	}

	d := make(map[string]int64)
	d["连接用时(ms)"] = (end - start)/1e6
	c.ReturnResponse(models.SUCCESS, d, true)
}

// 获取主机列表
func (c *HostController)GetHosts() {
	userId := c.LoginRequired(true)

	// 获取自己关注的主机信息
	myWatchs := make([]models.HostProfile, 0)
	c.o = orm.NewOrm()
	qs := c.o.QueryTable(new(models.HostWatch))
	var hws []*models.HostWatch
	_, _ = qs.Filter("user_id__exact", userId).All(&hws)
	for _, hw := range hws {
		h, _ := models.GetHostById(hw.HostId.Id)
		myWatchs = append(myWatchs, h.Host2Profile())
	}

	// 获取其他主机信息
	notWatchs := make([]models.HostProfile, 0)
	var hs []models.Host
	qs = c.o.QueryTable(new(models.Host))
	_, _ = qs.All(&hs)
	for _, h := range hs {
		// 筛选 host_watch 中 host_id=h.Id and user_id=userId
		cnt, _ := c.o.QueryTable(new(models.HostWatch)).Filter("host_id", h.Id).Filter("user_id", userId).Count()
		// 如果没有，说明没有关注这个主机，添加到返回值中
		if cnt == 0 {
			notWatchs = append(notWatchs, h.Host2Profile())
		}
	}

	d := make(map[string]interface{})
	d["my_watchs"] = myWatchs
	d["not_watchs"] = notWatchs
	c.ReturnResponse(0, d, true)
}

// TODO 开启主机监控

type SSHClientConfig struct {
	Host       string       //ip
	Port       int64        // 端口
	Username   string       //用户名
	Password   string       //密码
	Client	   *ssh.Client //ssh client
	LastResult string       //最近一次运行的结果
}

func (cliConf *SSHClientConfig)CreateClient() error {
	config := ssh.ClientConfig{
		User:              cliConf.Username,
		Auth:              []ssh.AuthMethod{ssh.Password(cliConf.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout:           120 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)
	cli, err := ssh.Dial("tcp", addr, &config)
	if  err != nil {
		logrus.Errorf("connection error: %v", err)
		return err
	}
	cliConf.Client = cli
	return nil
}

func (cliConf *SSHClientConfig) RunShell(shell string) (err error) {
	session, err := cliConf.Client.NewSession()
	if err != nil {
		return err
	}

	res, err := session.CombinedOutput(shell)
	if err != nil {
		return err
	}

	cliConf.LastResult = string(res)
	return nil
}
