package model

import (
	"github.com/cs-shuai/go-api-testing/common"
	"gopkg.in/check.v1"
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
	// fmt.Println("---------------" + fmt.Sprint(ob) + "---------------")
	ob.Value("msg").Equal("密码不正确")
}

func (l *Login) TestLoginSuccess(c *check.C) {
	// common.ParamByJson(l, common.RootPath+ "/json/login.json")
	l.Username = "admin"
	l.Password = "jcc2018"
	// fmt.Println("------ParamByJson---------" + fmt.Sprint(l) + "---------------")
	res := common.HttpPost(c, l)
	ob := res.JSON().Object()
	// fmt.Println("---------------" + fmt.Sprint(ob) + "---------------")
	ob.Value("token").NotNull()
}

func (l *Login) TearDownTest(_ *check.C) {
	_, err := l.Db.Exec("UPDATE user_token SET merchid=1  WHERE id=1")
	if err != nil {
		panic(err)
	}
}
