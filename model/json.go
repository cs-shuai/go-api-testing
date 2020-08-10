package model

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
	"jccAPITest/common"
	"jccAPITest/validation"
	"runtime"
	"strings"
	"sync"
)

const (
	GET  = "Get"
	POST = "Post"
)

var isLogin bool
var jsonTestingWait *sync.WaitGroup

/**
 * 自动测试
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type JsonTesting struct {
	RequestData     map[string]interface{} `json:"request_data"  mapstructure:"request_data" `
	RequestDataList []interface{}          `json:"request_data_list"  mapstructure:"request_data_list" `
	RequestUrl      string                 `json:"request_url"  mapstructure:"request_url"`
	RequestDataUrl  string                 `json:"request_data_url"  mapstructure:"request_data_url"`
	RequestMethod   string                 `json:"type"  mapstructure:"type"`
	TokenKey        string                 `json:"addr"  mapstructure:"addr"`
	TokenIsHeader   string                 `json:"is_header"  mapstructure:"is_header"`
	C               *check.C
	common.BaseJccAPITesting
}

func (jt *JsonTesting) UrlPath() (s string) { return jt.RequestUrl }

/**
 * Post参数处理
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) HandleParam(gat common.GoApiTesting) map[string]interface{} {
	common.ParamRandomByMap(&jt.RequestData)
	return jt.RequestData
}

/**
 * Get参数处理
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) HandleUrlCode(gat common.GoApiTesting) map[string]interface{} {
	return jt.RequestData
}

func (jt *JsonTesting) SetUpTest(c *check.C) {
	jt.Validation()
}
func (jt *JsonTesting) TearDownSuite(c *check.C) {
	jt.Validation()
}
func (jt *JsonTesting) SetUpSuite(c *check.C) {
	if !isLogin {
		jt.C = c
		// fmt.Println("-----------JsonTesting-SetUp---" + fmt.Sprint() + "---------------")
		login := new(Login)
		login.TestLoginSuccess(jt.C)
		pmToken := login.Response.JSON().Object().Raw()["token"].(string)
		jt.AddHeader("PmToken", pmToken)
		jt.AddParams("PmToken", pmToken)

		storeLogin := new(StoreLogin)
		storeLogin.TestLoginSuccess(jt.C)
		token := storeLogin.Response.JSON().Object().Raw()["token"].(string)
		jt.AddHeader("Token", token)
		jt.AddParams("Token", token)

		// jt.AddHeader("PmToken", "e611c721-a4d2-40cb-a8c3-6c8995e495b5")
		// jt.AddHeader("Token", "a2ae8a7e-b705-43ea-bd38-c91721aab743")
		// jt.AddHeader("MemberToken", "a2ae8a7e-b705-43ea-bd38-c91721aab743")
		isLogin = true
	}
	jt.Validation()
}
func (jt *JsonTesting) TearDownTest(c *check.C) {
	jt.Validation()
	jt.GetWaitGroup().Done()
}

func (jt *JsonTesting) SetUp()    {}
func (jt *JsonTesting) TearDown() {}
func (jt *JsonTesting) GetRouteDir() string {
	return viper.GetString("JSON_ROUTE_PATH")
}

/**
 * 添加头数据
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) AddHeader(key, value string) {
	common.AddHeaderGlobal(key, value)
}

/**
 * 添加参数数据
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) AddParams(key string, value interface{}) {
	common.AddParamsGlobal(key, value)
}

/**
 * 获取头部数据
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) GetHeader() map[string]string {
	return common.GetHeaderGlobal()
}

/**
 * 获取参数数据
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) GetParams() map[string]interface{} {
	return common.GetParamsGlobal()
}

/**
 * 初始化测试
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func (jt *JsonTesting) NewTesting(requestData interface{}) common.AutoTesting {
	newJsonTesting := new(JsonTesting)
	if err := mapstructure.Decode(requestData, newJsonTesting); err != nil {
		panic(err)
	}

	// 存在请求参数地址 获取参数并合并
	if newJsonTesting.RequestDataUrl != "" {
		paramArr := common.GetParamsByJsonFile(newJsonTesting.RequestDataUrl, viper.GetString("JSON_PATH"))
		newJsonTesting.RequestDataList = append(newJsonTesting.RequestDataList, paramArr...)
	}

	// 注册到测试中
	for _, test := range newJsonTesting.RequestDataList {
		jt.GetWaitGroup().Add(1)
		jTemp := *newJsonTesting
		jTemp.RequestData = test.(map[string]interface{})
		check.Suite(&jTemp)
	}

	return newJsonTesting
}

/**
 * 测试执行方法 [通过json生成不同数据的对象 最后执行的方法]
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func (jt *JsonTesting) TestRun(c *check.C) common.AutoTesting {
	jt.C = c
	jt.Request()

	return jt
}

/**
 * 请求
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func (jt *JsonTesting) Request() common.AutoTesting {
	switch jt.RequestMethod {
	case GET:
		common.HttpGet(jt.C, jt)
	case POST:
		common.HttpPost(jt.C, jt)
	default:
		common.HttpPost(jt.C, jt)
	}

	return jt
}

/**
 * 回应校验
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func (jt *JsonTesting) ResponseCheck() common.AutoTesting {
	// fmt.Println("---------jt.Response------" + fmt.Sprint(jt.Response) + "---------------")

	responseKey := common.ResponseKey
	response := common.Response
	equalValue := "成功"
	if value, ok := jt.RequestData[response]; ok {
		equalValue = value.(string)
	}

	jt.Response.JSON().Object().Value(responseKey).Equal(equalValue)
	return jt
}

/**
 * 获取WaitGroup
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) GetWaitGroup() *sync.WaitGroup {
	return jsonTestingWait
}

/**
 * 初始化WaitGroup
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) InitWaitGroup() {
	jsonTestingWait = new(sync.WaitGroup)
}

/**
 * 调取验证类
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) Validation() {
	funcName := getFuncName(2)
	for _, v := range validation.CheckList {
		if v.GetRunFunc() == funcName {
			if value, ok := jt.RequestData[v.GetJsonKey()]; ok {
				v.SetJsonValue(value)
				v.Run(jt.Response, &jt.RequestData)
			}
			// fmt.Println("---------------" + fmt.Sprint(jt.RequestData) + "---------------")
		}
	}
}

/**
 * 获取方法名称
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func getFuncName(skip int) string {
	// 获取上一个执行方法位置
	pc, _, _, _ := runtime.Caller(skip)
	// 获取方法名
	fn := runtime.FuncForPC(pc).Name()
	mn := strings.Split(fn, ".")

	return mn[len(mn)-1]
}
