package model

import (
	"github.com/cs-shuai/go-api-test/common"
	"gopkg.in/check.v1"
)

type StoreLogin struct {
	Username string `json:"username"`
	Password string `json:"pwd"`
	common.BaseJccAPITesting
}

func init() {
	common.RegisterCheck(new(StoreLogin))
}

func (l *StoreLogin) UrlPath() string {
	return "store/account/user_login"
}

func (l *StoreLogin) TestLoginError(c *check.C) {
	l.Username = "admin"
	l.Password = "23232"
	res := common.HttpPost(c, l)
	ob := res.JSON().Object()
	// fmt.Println("---------------" + fmt.Sprint(ob) + "---------------")
	ob.Value("msg").Equal("登录失败")
}

func (l *StoreLogin) TestLoginSuccess(c *check.C) {
	// common.ParamByJson(l, common.RootPath+ "/json/login.json")
	l.Username = "大成测试"
	l.Password = "123"
	// fmt.Println("------ParamByJson---------" + fmt.Sprint(l) + "---------------")
	res := common.HttpPost(c, l)
	ob := res.JSON().Object()
	// fmt.Println("---------------" + fmt.Sprint(ob) + "---------------")
	ob.Value("token").NotNull()
}
