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
	userId := c.LoginRequired()
	userObj, err := models.GetUserById(userId)
	if err != nil {
		c.ReturnResponse(models.AUTH_ERROR, nil, true)
	}

	var req models.AddHostReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.ReturnResponse(models.REQUEST_ERROR, nil, true)
	}

	c.o = orm.NewOrm()
	hostObj := models.Host{
		UserId:        userObj,
		CreateTime:    time.Time{},
		Ip:            req.Ip,
		Name:          req.Name,
		Description:   req.Description,
		LoginName:     req.LoginName,
		LoginPwd:      req.LoginPwd,
	}

	if _, err = models.AddHost(&hostObj); err != nil {
		logrus.Warningf("User:%v 添加主机失败，Request：%v，错误信息：%v", userId, req, err)
		c.ReturnResponse(models.SERVER_ERROR, nil, true)
	}

	// TODO 开一个协程去获取主机 base info
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

// 测试主机连接
func (c *HostController) HostConnectionTest()  {
	c.LoginRequired()

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
