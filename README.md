# 接口测试

## 执行  
* 根目录执行
```
go test
```

## 创建新项目
1. 获取测试框架包
```
go get github.com/cs-shuai/go-api-testing
```

2. 创建 main_test文件 选择注册对象
```
package go_api_testing_test

import (
	"github.com/cs-shuai/go-api-testing/common"
	"github.com/cs-shuai/go-api-testing/model"
	"testing"
)

func Test(t *testing.T) {
    // 选择注册测试对象
	common.AutoTestRun(t, new(model.JsonTesting), new(model.JsonGroupTesting))
}

```
3. 创建配置文件 conf/ config.toml
```
Host = "https://XXXXXXXXXXX/"
TOKEN_KEY = "XXXXXX"
JSON_PATH = "json/"
JSON_ROUTE_PATH = "route/"
JSON_GROUP_ROUTE_PATH = "group_route/"
SQLCONN = "user:pass@tcp(127.0.0.1:3306)/database"
```
4. 创建对应的json目录和json文件
5. 项目根目录 执行
```
go test
```


## JSON文件规则
* json文件中为数组对象
* 请求参数可在对应文件中编写 同样为数组格式
* 请求参数(request_data) 中可随机生成参数 
    * auto 随机字符串
    * auto_int 随机数组
```
[
  {
    "request_url": "manager/pm_user/user_login",
    "request_data_url": "login.json",
    "request_data" :[],
    "type": "Post",
    "addr": "PmToken"
  },
  {
    "request_url": "manager/pm_member/select_members",
    "request_data_url": "select_member.json",
    "type": "Get",
    "addr": "PmToken"
  }
]

```

## 编写测试用例
* <a href="#结构体创建">结构体 </a>
* <a href="#Json创建">json </a>
* <a href="#GroupJson创建">GroupJson </a>


## 结构体创建
<a name="1">结构体创建</a>
1. 实现接口 AutoTesting
2. 继承BaseJccAPITesting  方便实现些可必须的方法
3. 在实现结构体下创建 前缀为 `Test` 方法
```
package model

import (
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

```
## Json创建
<a name="2">Json创建</a>
1. 在配置的指定route目录中创建json文件 编写对应的请求地址和请求参数  自动执行
```
[
  {
    "request_url": "manager/pm_user/user_login",  // 请求地址
    "request_data_url": "login.json",             // 请求数据地址 (目录在配置中配置)
    "type": "Post",                               // 请求方式 
    "addr": "PmToken"                             // 权限验证 
    "request_data": [{
                "username": "丛力强",
                "pwd": "123",
                "response" : "操作成功"
              },
              {
                "username": "丛力强1",
                "pwd": "1232",
                "response" : "用户不存在"
              }
            ]},                                     // 请求参数  可写在对应文件中  2数据会合并
  {
    "request_url": "manager/pm_member/select_members",
    "request_data_url": "select_member.json",
    "type": "Get",
    "addr": "PmToken"
  }
]

```

## GroupJson创建
<a name="1">GroupJson创建</a>
1. 基础配置与上述一致 区别在于  一个测试用例必须写在同一个文件中 
> 一个用例可以请求多个接口  前置接口为后续接口获取参数
```
[
  {
    "request_url": "manager/pm_user/user_login",
    "request_data": [
      {
        "username": "丛力强",
        "pwd": "123",
        "response" : "操作成功"
      }
    ],
    "type": "Post",
    "addr": "PmToken"
  },
  {
    "request_url": "manager/pm_member/condition_list",
    "request_data": [
      {
        "response" : "操作成功",
        "before": [{
          "url" : "manager/pm_user/user_login",
          "before_key" : "token",
          "key" : "PmToken",
          "is_header" : true
        }]
      }
    ],
    "type": "Get",
    "addr": "PmToken"
  },
  {
    "request_url": "manager/pm_member/pm_add_member_tag",
    "request_data": [
      {
        "response" : "操作成功",
        "merch_id": "undefined",
        "tag": "auto",
        "description": "",
        "member_ids": "",
        "recommend": 0,
        "and_or": 0,
        "tag_type": 1,
        "condition_options_str": [
          {
            "condition_id": "1",
            "condition_value": []
          },
          {
            "condition_id": "5",
            "condition_value": [
              9
            ]
          }
        ],
        "before": [{
          "url" : "manager/pm_user/user_login",
          "before_key" : "token",
          "key" : "PmToken",
          "is_header" : true
        }, {
            "url" : "manager/pm_member/condition_list",
            "before_key" : "data.random.id",
            "key" : "tag_type"
          }]
      }
    ],
    "type": "Post",
    "addr": "PmToken"
  }
]

```

## 自定义验证标签
> 指定的json键值 执行 对应时机的对应回调方法
> 只需实现 BaseJsonKeyExpandInterface
```
package json_key

import "github.com/spf13/viper"

var ResponseMessage string
func init() {
	RegisterJsonKey(new(Response))
}

/**
 * 返回匹配
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type Response struct {
	Key   string
	Value interface{}
}

func (r *Response) GetJsonKey() string {
	r.Key = "response"
	return r.Key
}
func (r *Response) SetJsonValue(value interface{}) {
	r.Value = value
}

func (r *Response) GetJsonValue() interface{} {
	return r.Value
}

func (r *Response) TearDownRun(params *J) {
	if ResponseMessage == "" {
		ResponseMessage = viper.GetString("RESPONSE_MESSAGE")
	}
	params.Response.JSON().Object().Value(ResponseMessage).Equal(r.Value)
}

func (r *Response) SetUpRun(params *J) {

}

```

## Json值的拓展
```
package json_value

import (
	"math/rand"
	"time"
)

func init() {
	RegisterJsonValue(new(Auto))
}
type Auto struct {
}

func (a *Auto) GetJsonValue() string {
	return "auto"
}

func (a *Auto) Run() interface{} {
	return 	randomString(8)
}


/**
 * 随机字符串
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func randomString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
```

## 数据库操作
> 通过设置配置文件 `SQLCONN` 连接数据库
> 在 BaseJccAPITesting 结构体中存在 DB 属性 通过方法初始化
```
func (t *BaseJccAPITesting) Initialization() {
	db, err := sql.Open("mysql", viper.GetString("SQLCONN"))
	if err != nil {
		panic(err)
	}

	t.Db = db
}
```

## 测试UML图
![img](./common/go_api_testing.png)