package model

import (
	"fmt"
	"gopkg.in/check.v1"
	"jccAPITest/common"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"pwd"`
	common.BaseJccAPITesting
}

func init() {
	common.RegisterCheck(new(Login))
}

func (l *Login) UrlPath() string {
	return "manager/pm_user/user_login"
}

func (l *Login) TestLoginError(c *check.C) {
	l.Username = "admin"
	l.Password = "23232"
	res := common.HttpPost(c, l)
	ob := res.JSON().Object()
	fmt.Println("---------------" + fmt.Sprint(ob) + "---------------")
	ob.Value("msg").Equal("密码不正确")
}

func (l *Login) TestLoginSuccess(c *check.C) {
	// l.Username = "admin"
	// l.Password = "jcc2018"
	common.ParamByJson(l, "../params/login.json")
	fmt.Println("------ParamByJson---------" + fmt.Sprint(l) + "---------------")
	res := common.HttpPost(c, l)
	ob := res.JSON().Object()
	fmt.Println("---------------" + fmt.Sprint(ob) + "---------------")
	ob.Value("token").NotNull()
}
