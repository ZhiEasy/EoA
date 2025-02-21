package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
	BaseController
}

func (c *UserController) GetYuQueOAuthPath() {
	state := c.GetString("state")
	if state == "" {
		state = "default"
	}
	/*
		https://www.yuque.com/oauth2/authorize?
		client_id=FCEGPMmDcnjwDKJsTfoV
		&scope=group:read
		&redirect_uri=http://127.0.0.1:10240/user/oauth/yuque
		&state=123456
		&response_type=code
	*/
	clientId := beego.AppConfig.String("YuQue::ClientID")
	redirectUrl := beego.AppConfig.String("YuQue::RedirectURI")
	s := fmt.Sprintf("https://www.yuque.com/oauth2/authorize?client_id=%s&state=%s&response_type=code&scope=group:read&redirect_uri=%s", clientId, state, redirectUrl)
	c.ReturnResponse(models.SUCCESS, s, true)
}

/*
语雀 OAuth
点击 https://www.yuque.com/oauth2/authorize?client_id=FCEGPMmDcnjwDKJsTfoV&scope=group:read&redirect_uri=http://127.0.0.1:10240/user/oauth&state=123456&response_type=code
授权后的回调接口
*/
func (c *UserController) YuQueOAuthRedirect() {
	// 解析参数
	code := c.GetString("code")
	state := c.GetString("state")

	authRedirectURL := beego.AppConfig.String("YuQue::AuthRedirectTo")
	retUrlValue := url.Values{}
	retUrlValue.Add("state", state)

	// 根据 code 换取用户 token
	token, err := GetUserToken(code)
	if err != nil {
		c.Redirect(authRedirectURL+"?"+retUrlValue.Encode(), 302)
	}
	// 根据 token 换取用户信息
	userInfo, err := GetUserInfo(token)
	if err != nil {
		c.Redirect(authRedirectURL+"?"+retUrlValue.Encode(), 302)
	}
	// 检查用户是否在组织中
	ok, err := CheckUserInGroup(userInfo)
	if err != nil {
		c.Redirect(authRedirectURL+"?"+retUrlValue.Encode(), 302)
	}

	// 用户不在组织中
	if !ok {
		c.Redirect(authRedirectURL+"?"+retUrlValue.Encode(), 302)
	}

	c.o = orm.NewOrm()
	var user models.User
	// 根据语雀返回的用户id判断用户是否已经创建过了
	qs := c.o.QueryTable(user) // user相当于"user"，表示查user表
	err = qs.Filter("yuque_id__exact", userInfo.Data.ID).One(&user)
	user.YuqueId = userInfo.Data.ID

	var id int64
	// 没有找到，新用户
	if err != nil {
		// 保存语雀返回的用户信息到 yuque_info
		b, _ := json.Marshal(userInfo)
		user.YuqueInfo = string(b)
		// 添加用户
		id, _ = models.AddUser(&user)
		c.SetSession("user_id", id)
		c.Redirect(authRedirectURL+"?"+retUrlValue.Encode(), 302)
	}

	// 已经完善了信息
	id = int64(user.Id)
	c.SetSession("user_id", user.Id)
	c.Redirect(authRedirectURL+"?"+retUrlValue.Encode(), 302)
}

/*
Github OAuth
*/
func (c *UserController) GithubOAuthRedirect() {
}

// 用户登录
func (c *UserController) UserLogin() {
	var req models.UserLoginReq
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.ReturnResponse(models.REQUEST_ERROR, nil, true)
	}

	hash := md5.New()
	hash.Write([]byte(req.Pwd))
	md5Pwd := hex.EncodeToString(hash.Sum(nil))

	var userObj models.User
	c.o = orm.NewOrm()
	qs := c.o.QueryTable(new(models.User))
	err = qs.Filter("email", req.Email).Filter("pwd", md5Pwd).One(&userObj)
	if err != nil {
		c.ReturnResponse(models.PWD_ERROR, nil, true)
	}

	c.SetSession("user_id", userObj.Id)
	c.ReturnResponse(models.SUCCESS, userObj.User2UserProfile(), true)
}

// 退出登录
func (c *UserController) UserLogout() {
	_ = c.LoginRequired(false)
	c.DelSession("user_id")
	c.ReturnResponse(models.SUCCESS, nil, true)
}

// 用户完善信息接口
func (c *UserController) UpdateUserInfo() {
	userId := c.LoginRequired(false)

	var req models.UpdateUserInfoReq
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.ReturnResponse(models.REQUEST_ERROR, nil, true)
	}

	user, err := models.GetUserById(userId)
	if err != nil {
		c.ReturnResponse(models.AUTH_ERROR, nil, true)
		return // 让下面的 user 不报警告
	}

	if req.Pwd != req.CPwd {
		c.ReturnResponse(models.CPWD_ERROR, nil, true)
	}

	user.Name = req.Name
	user.Email = req.Email
	_, _ = md5.New().Write([]byte(req.Pwd))
	hash := md5.New()
	hash.Write([]byte(req.Pwd))
	user.Pwd = hex.EncodeToString(hash.Sum(nil))

	// 更新用户信息
	if err = models.UpdateUserById(user); err != nil {
		c.ReturnResponse(models.SERVER_ERROR, nil, true)
	}

	c.ReturnResponse(models.SUCCESS, nil, true)
}

// 获取当前登录用户信息接口
func (c *UserController) GetUserInfo() {
	userId := c.LoginRequired(false)

	user, err := models.GetUserById(userId)
	if err != nil {
		c.ReturnResponse(models.AUTH_ERROR, nil, true)
		return
	}

	var yuque models.YuQueUserInfo
	_ = json.Unmarshal([]byte(user.YuqueInfo), &yuque)
	userInfo := models.UserProfile{
		Id:         user.Id,
		CreateTime: user.CreateTime,
		Name:       user.Name,
		Email:      user.Email,
		AvatarUrl:  yuque.Data.AvatarURL,
	}

	// 判断这个已经授权过的用户是否完善了信息
	if user.Pwd == "" || user.Name == "" || user.Email == "" {
		// 如果没有完善信息
		c.ReturnResponse(models.NEED_UPDATE_INFO, userInfo, true)
	}
	c.ReturnResponse(models.SUCCESS, userInfo, true)
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
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
		GrantType    string `json:"grant_type"`
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
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
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
