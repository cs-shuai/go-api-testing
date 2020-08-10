package model

import (
	"github.com/gavv/httpexpect"
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

var jsonTestingWait sync.WaitGroup

type JsonTesting struct {
	RequestData      map[string]interface{} `json:"request_data"  mapstructure:"request_data" `
	RequestDataList  []interface{}          `json:"request_data_list"  mapstructure:"request_data_list" `
	RequestUrl       string                 `json:"request_url"  mapstructure:"request_url"`
	RequestDataUrl   string                 `json:"request_data_url"  mapstructure:"request_data_url"`
	RequestMethod    string                 `json:"type"  mapstructure:"type"`
	TokenKey         string                 `json:"addr"  mapstructure:"addr"`
	TokenIsHeader    string                 `json:"is_header"  mapstructure:"is_header"`
	ResponseToParams []*ResponseToParam     `json:"response_param"  mapstructure:"response_param"`
	TestResultList   []*TestResult
	C                *check.C
	common.BaseJccAPITesting
}

type TestResult struct {
	RequestData      map[string]interface{} `json:"request_data"  mapstructure:"request_data" `
	Response         *httpexpect.Response   `json:"response"  mapstructure:"response"`
	ResponseToParams []*ResponseToParam     `json:"response_param"  mapstructure:"response_param"`
}

type ResponseToParam struct {
	ResponseKey string
	ParamKey    string
}

func (jt *JsonTesting) UrlPath() (s string) { return jt.RequestUrl }
func (jt *JsonTesting) HandleParam(gat common.GoApiTesting) map[string]interface{} {
	return jt.RequestData
}
func (jt *JsonTesting) HandleUrlCode(gat common.GoApiTesting) map[string]interface{} {
	return jt.RequestData
}

var isLogin bool

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
		token := login.Response.JSON().Object().Raw()["token"].(string)
		jt.AddHeader(jt.TokenKey, token)
		jt.AddParams(jt.TokenKey, token)

		isLogin = true
	}
	jt.Validation()
}
func (jt *JsonTesting) SetUp() {}
func (jt *JsonTesting) TearDown() {
	// fmt.Println("-----------JsonTesting-TearDown---" + fmt.Sprint() + "---------------")
}
func (jt *JsonTesting) AddHeader(key, value string) {
	common.AddHeaderGlobal(key, value)
}

func (jt *JsonTesting) AddParams(key string, value interface{}) {
	common.AddParamsGlobal(key, value)
}

func (jt *JsonTesting) GetHeader() map[string]string {
	return common.GetHeaderGlobal()
}

func (jt *JsonTesting) GetParams() map[string]interface{} {
	return common.GetParamsGlobal()
}

/**
 * 初始化测试
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func (jt *JsonTesting) NewTesting(requestData interface{}) common.AutoTesting {
	// fmt.Println("---------------" + fmt.Sprint(jt) + "---------------")
	if err := mapstructure.Decode(requestData, jt); err != nil {
		panic(err)
	}

	if jt.RequestDataUrl != "" {
		paramArr := common.GetParamsByJsonFile(jt.RequestDataUrl, viper.GetString("JSON_PATH"))
		jt.RequestDataList = append(jt.RequestDataList, paramArr...)
	}
	// fmt.Println("---------testRequest------" + fmt.Sprint(jt) + "---------------")
	for _, test := range jt.RequestDataList {
		jt.GetWaitGroup().Add(1)
		jt.RequestData = test.(map[string]interface{})
		check.Suite(jt)
	}

	return jt
}

/**
 * 测试执行方法 [通过json生成不同数据的对象 最后执行的方法]
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func (jt *JsonTesting) TestRun(c *check.C) common.AutoTesting {
	jt.C = c
	// fmt.Println("---------------" + fmt.Sprint(test) + "---------------")
	jt.Request()

	return jt
}

/**
 * 请求
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func (jt *JsonTesting) Request() common.AutoTesting {
	// fmt.Println("--------Request-------" + fmt.Sprint(jt) + "---------------")
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
	jt.Response.JSON().Object().Value(responseKey).Equal("111")
	return jt
}

func (jt *JsonTesting) GetWaitGroup() *sync.WaitGroup {
	return &jsonTestingWait
}

func (jt *JsonTesting) TearDownTest(c *check.C) {
	jt.Validation()
	jt.GetWaitGroup().Done()

}

func (jt *JsonTesting) Validation() {
	funcName := getFuncName(2)
	for _, v := range validation.CheckList {
		if v.GetRunFunc() == funcName {
			// fmt.Println("---------------" + fmt.Sprint(jt.RequestData) + "---------------")
			if value, ok := jt.RequestData[v.GetJsonKey()]; ok {
				// fmt.Println("---------------" + fmt.Sprint(value) + "---------------")

				v.SetJsonValue(value)
				v.Run(jt.Response, &jt.RequestData)
			}
			// fmt.Println("---------------" + fmt.Sprint(jt.RequestData) + "---------------")
		}
	}
}

func getFuncName(skip int) string {
	// 获取上一个执行方法位置
	pc, _, _, _ := runtime.Caller(skip)
	// 获取方法名
	fn := runtime.FuncForPC(pc).Name()
	mn := strings.Split(fn, ".")

	return mn[len(mn)-1]
}
