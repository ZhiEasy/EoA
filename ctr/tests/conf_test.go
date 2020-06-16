package test

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"testing"
)

func TestConfigToken(t *testing.T) {
	fmt.Println(beego.BConfig.AppName)

	conf, err := config.NewConfig("json", "../conf/ctr.json")
	//fmt.Println(conf)
	fmt.Println(conf, err)
	val := conf.String("dev::yuquetoken")
	fmt.Println(val)
}