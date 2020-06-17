package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/ahojcn/EoA/ctr/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"net/url"
)

// UserController operations for User
type UserController struct {
	beego.Controller
	o orm.Ormer
}
/*
https://www.yuque.com/oauth2/authorize
?client_id=FCEGPMmDcnjwDKJsTfoV
&scope=group:read
&redirect_uri=http://127.0.0.1:10240/user/oauth
&state=123456
&response_type=code
*/
func (c *UserController)OAuth() {
	// 解析参数
	code := c.GetString("code")
	state := c.GetString("state")

	authRedirectURL := beego.AppConfig.String("YuQue::AuthRedirectTo")
	retUrlValue := url.Values{}
	retUrlValue.Add("state", state)

	// 根据 code 换取用户 token
	token, err := GetUserToken(code)
	if err != nil {
		retUrlValue.Add("status", "-1")
		c.Redirect(authRedirectURL + "?" + retUrlValue.Encode(), 302)
	}
	// 根据 token 换取用户信息
	userInfo, err := GetUserInfo(token)
	if err != nil {
		retUrlValue.Add("status", "-1")
		c.Redirect(authRedirectURL + "?" + retUrlValue.Encode(), 302)
	}
	// 检查用户是否在组织中
	ok, err := CheckUserInGroup(userInfo)
	if err != nil {
		retUrlValue.Add("status", "-1")
		c.Redirect(authRedirectURL + "?" + retUrlValue.Encode(), 302)
	}

	// 用户不在组织中
	if !ok {
		retUrlValue.Add("status", "-1")
		c.Redirect(authRedirectURL + "?" + retUrlValue.Encode(), 302)
	}

	c.o = orm.NewOrm()
	var user models.User
	// 判断用户是否已经创建过了
	qs := c.o.QueryTable(user)  // user相当于"user"，表示查user表
	err = qs.Filter("yuque_token__exact", token).One(&user)

	var id int64
	if err != nil {
		// 没有找到，新用户
		id, _ = models.AddUser(&user)
	} else {
		// 找到了，老用户（已经yuque授权过）
		id = int64(user.Id)
	}

	retUrlValue.Add("status", "-1")
	retUrlValue.Add("id", string(id))
	c.Redirect(authRedirectURL + "?" + retUrlValue.Encode(), 302)
}

// 获取组织的用户，检查用户是否在组织中
func CheckUserInGroup(yuqueUserInfo *models.YuQueUserInfo) (ok bool, err error) {
	// 获取组织信息
	httpCli := http.Client{}
	request, err := http.NewRequest("GET", "https://www.yuque.com/api/v2/groups/1167287/users", nil)
	if err != nil {
		logrus.Fatalf("获取组织信息失败 新建请求失败 %v", err)
		return false, errors.New("检查用户权限失败")
	}
	request.Header.Add("X-Auth-Token", beego.AppConfig.String("YuQue::Token"))
	resp, err := httpCli.Do(request)
	if err != nil {
		logrus.Fatalf("获取组织信息失败 获取组织信息失败 %v", err)
		return false, errors.New("检查用户权限失败")
	}
	respBytes, err := ioutil.ReadAll(resp.Body)

	var usersInGroup models.YuQueGroupUsers
	_ = json.Unmarshal(respBytes, &usersInGroup)

	// 检查用户是否在组织中
	ok = false
	for _, user := range usersInGroup.Data {
		if user.UserID == yuqueUserInfo.Data.ID {
			ok = true
		}
	}

	return ok, nil
}

// 根据code换取用户token
func GetUserToken(code string) (token string, err error) {
	// 获取用户 token
	clientID := beego.AppConfig.String("YuQue::ClientID")
	clientSecret := beego.AppConfig.String("YuQue::ClientSecret")
	grantType := "authorization_code"
	type tmp struct {
		ClientID string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code string `json:"code"`
		GrantType string `json:"grant_type"`
	}
	req := tmp{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Code:         code,
		GrantType:    grantType,
	}
	byteData, err := json.Marshal(&req)
	if err != nil {
		logrus.Fatalf("换取token失败 解析参数错误 %v", err)
		return "", errors.New("获取用户信息失败")
	}
	reader := bytes.NewReader(byteData)
	request, err := http.NewRequest("POST", "https://www.yuque.com/oauth2/token", reader)
	if err != nil {
		logrus.Fatalf("换取token失败 创建新请求失败 %v", err)
		return "", errors.New("获取用户信息失败")
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpCli := http.Client{}
	resp, err := httpCli.Do(request)
	if err != nil {
		logrus.Fatalf("换取token失败 请求失败 %v", err)
		return "", errors.New("获取用户信息失败")
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("换取token失败 读取请求的响应失败 %v", err)
		return "", errors.New("获取用户信息失败")
	}
	type respToken struct {
		AccessToken string `json:"access_token"`
		TokenType string `json:"token_type"`
		Scope string `json:"scope"`
	}
	var r respToken
	err = json.Unmarshal(respBytes, &r)
	if err != nil {
		logrus.Fatalf("换取token失败 解析请求返回值失败 %v", err)
		return "", errors.New("获取用户信息失败")
	}

	return r.AccessToken, nil
}

// 根据token换取用户信息
func GetUserInfo(token string) (yuqueUserInfo *models.YuQueUserInfo, err error) {
	httpCli := http.Client{}
	// 根据 token 获取用户信息
	request, err := http.NewRequest("GET", "https://www.yuque.com/api/v2/user", nil)
	if err != nil {
		logrus.Fatalf("根据token换取用户信息失败 新建获取token的请求失败 %v", err)
		return nil, errors.New("获取用户信息失败")
	}
	request.Header.Add("X-Auth-Token", token)
	resp, err := httpCli.Do(request)
	if err != nil {
		logrus.Fatalf("根据token换取用户信息失败 根据token获取用户信息失败 %v", err)
		return nil, errors.New("获取用户信息失败")
	}
	respBytes, err := ioutil.ReadAll(resp.Body)

	var userInfo models.YuQueUserInfo
	err = json.Unmarshal(respBytes, &userInfo)
	if err != nil {
		logrus.Fatalf("根据token换取用户信息失败 解析获取到的用户信息错误 %v", err)
		return nil, errors.New("获取用户信息失败")
	}

	return &userInfo, nil
}