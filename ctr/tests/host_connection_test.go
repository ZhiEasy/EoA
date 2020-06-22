package test

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net"
	"testing"
	"time"
)

//连接的配置
type ClientConfig struct {
	Host       string       //ip
	Port       int64        // 端口
	Username   string       //用户名
	Password   string       //密码
	Client	   *ssh.Client //ssh client
	LastResult string       //最近一次运行的结果
}

var clientConfig ClientConfig

func TestHostConnection(t *testing.T) {
	clientConfig.Host = "47.103.14.73"
	clientConfig.Port = 22
	clientConfig.Username = "root"
	clientConfig.Password = "ndl.04551"

	config := ssh.ClientConfig{
		User:              clientConfig.Username,
		Auth:              []ssh.AuthMethod{ssh.Password(clientConfig.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout:           10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", clientConfig.Host, clientConfig.Port)
	client, err := ssh.Dial("tcp", addr, &config)
	if  err != nil {
		logrus.Errorf("connection error: %v", err)
	}

	logrus.Errorln(client.User(), string(client.ClientVersion()), client.LocalAddr().String(), client.RemoteAddr().String())

	session ,err :=client.NewSession()
	output, err := session.CombinedOutput("echo ok")
	fmt.Println(string(output) == "ok\n")
}
