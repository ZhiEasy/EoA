package controllers

import (
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

	hostObj := models.Host{
		UserId:        userObj,
		CreateTime:    time.Time{},
		Ip:            req.Ip,
		Name:          req.Name,
		Description:   req.Description,
		LoginName:     req.LoginName,
		LoginPwd:      req.LoginPwd,
	}
	//hostId, err := models.AddHost(&hostObj)
	_, err = models.AddHost(&hostObj)
	if err != nil {
		logrus.Warningf("User:%v 添加主机失败，Request：%v，错误信息：%v", userId, req, err)
		c.ReturnResponse(models.SERVER_ERROR, nil, true)
	}

	// TODO 开一个协程去获取主机 base info
	//go func() {
	//	host, _ := models.GetHostById(int(hostId))
	//}()

	err = AddHostWatch(userObj.Id, hostObj.Id)
	if err != nil {
		c.ReturnResponse(models.HOST_REWATCH, nil, true)
	}

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

// GET 获取主机列表
func (c *HostController)GetHosts() {
	userId := c.LoginRequired(true)

	// 获取自己关注的主机信息
	var myWatchs []models.HostProfile
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
	var n []*models.HostWatch
	_, _ = qs.Exclude("user_id__exact", userId).All(&n)
	for _, hw := range n {
		h, _ := models.GetHostById(hw.HostId.Id)
		notWatchs = append(notWatchs, h.Host2Profile())
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
		Timeout:           10 * time.Second,
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
