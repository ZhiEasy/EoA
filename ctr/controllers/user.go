package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/ahojcn/EoA/ctr/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"time"
)

type UserController struct {
	beego.Controller
	o orm.Ormer
	resp models.Response
}

/*
https://www.yuque.com/oauth2/authorize?client_id=FCEGPMmDcnjwDKJsTfoV&scope=group:read&redirect_uri=http://127.0.0.1:10240/user/oauth&state=123456&response_type=code
*/
func (c *UserController)OAuth() {
	// 根据 code 换取用户 token
	code := c.GetString("code")
	state := c.GetString("state")
	logrus.Info(state)
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
	logrus.Infoln(req)
	byteData, _ := json.Marshal(&req)
	reader := bytes.NewReader(byteData)
	request, err := http.NewRequest("POST", "https://www.yuque.com/oauth2/token", reader)
	if err != nil {
		logrus.Fatalf("http.NewRequest %v", err)
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logrus.Fatalf("client.Do %v", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatalf("ioutil.ReadAll %v", err)
	}
	type respToken struct {
		AccessToken string `json:"access_token"`
		TokenType string `json:"token_type"`
		Scope string `json:"scope"`
	}
	var r respToken
	logrus.Info(json.Unmarshal(respBytes, &r))

	// 根据 token 获取用户信息
	request, err = http.NewRequest("GET", "https://www.yuque.com/api/v2/user", nil)
	if err != nil {
		logrus.Fatalf("获取用户信息失败 %v", err)
	}
	request.Header.Add("X-Auth-Token", r.AccessToken)
	resp, err = client.Do(request)
	if err != nil {
		logrus.Fatalf("client.Do %v", err)
	}
	respBytes, err = ioutil.ReadAll(resp.Body)

	type yuqueUserInfo struct {
		Data struct {
			ID int `json:"id"`
			Type string `json:"type"`
			SpaceID int `json:"space_id"`
			AccountID int `json:"account_id"`
			Login string `json:"login"`
			Name string `json:"name"`
			AvatarURL string `json:"avatar_url"`
			BooksCount int `json:"books_count"`
			PublicBooksCount int `json:"public_books_count"`
			FollowersCount int `json:"followers_count"`
			FollowingCount int `json:"following_count"`
			Public int `json:"public"`
			Description interface{} `json:"description"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Serializer string `json:"_serializer"`
		} `json:"data"`
	}
	var userInfo yuqueUserInfo
	_ = json.Unmarshal(respBytes, &userInfo)

	// 获取组织信息
	request, err = http.NewRequest("GET", "https://www.yuque.com/api/v2/groups/1167287/users", nil)
	if err != nil {
		logrus.Fatalf("获取组织信息失败 %v", err)
	}
	request.Header.Add("X-Auth-Token", beego.AppConfig.String("YuQue::Token"))
	resp, err = client.Do(request)
	if err != nil {
		logrus.Fatalf("获取组织信息失败 %v", err)
	}
	respBytes, err = ioutil.ReadAll(resp.Body)
	type yuqueGroupUsers struct {
		Data []struct {
			ID int `json:"id"`
			GroupID int `json:"group_id"`
			UserID int `json:"user_id"`
			User struct {
				ID int `json:"id"`
				Type string `json:"type"`
				Login string `json:"login"`
				Name string `json:"name"`
				Description interface{} `json:"description"`
				AvatarURL string `json:"avatar_url"`
				FollowersCount int `json:"followers_count"`
				FollowingCount int `json:"following_count"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
				Serializer string `json:"_serializer"`
			} `json:"user"`
			Role int `json:"role"`
			Status int `json:"status"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Group interface{} `json:"group"`
			Serializer string `json:"_serializer"`
		} `json:"data"`
	}
	var usersInGroup yuqueGroupUsers
	_ = json.Unmarshal(respBytes, &usersInGroup)

	ok := false
	for _, user := range usersInGroup.Data {
		if user.UserID == userInfo.Data.ID {
			ok = true
		}
	}

	//if ok {
	//	c.Redirect("http://www.baidu.com", 302)
	//} else {
	//	c.Redirect("http://google.com", 302)
	//}
	c.o = orm.NewOrm()
	logrus.Warnln(string(time.Now().Unix()))
	user := models.User{
		CreateTime: string(time.Now().Unix()),
		Name:       "",
		Email:      "",
		Pwd:        "",
		YuqueToken: r.AccessToken,
	}
	//user_id, _ := models.AddUser(&user)
	id, _ := c.o.Insert(&user)

	//u, _ := models.GetUserById(int(id))
	//logrus.Warnln(u.CreateTime)

	logrus.Infoln(ok)
	c.Data["json"] = id
	c.ServeJSON()
}