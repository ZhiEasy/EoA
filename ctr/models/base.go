package models

import (
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func Init() {
	logrus.Warnln(beego.AppConfig.String("sqlconn"));
	err := orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("sqlconn"))
	if err != nil {
		logrus.Panicf("注册数据库失败 %v", err)
	}
}

// 语雀用户信息
type YuQueUserInfo struct {
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

// 语雀用户组信息
type YuQueGroupUsers struct {
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

// 用户状态信息
type UserStatus struct {
	UID int64 `json:"uid"`
	IsLogin bool `json:"is_login"`
	Status int `json:"status"`
}