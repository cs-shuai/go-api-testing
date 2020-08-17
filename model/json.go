package model

import (
	"crypto/sha1"
	"fmt"
	"github.com/cs-shuai/go-api-testing/common"
	"github.com/cs-shuai/go-api-testing/json_key"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
	"runtime"
	"strings"
)

const (
	GET  = "Get"
	POST = "Post"
)

/**
 * 自动测试
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type JsonTesting struct {
	RequestData     map[string]interface{} `json:"request_data" `
	Header          map[string]string      `json:"header"  mapstructure:"header" `
	RequestDataList []interface{}          `json:"request_data_list"  mapstructure:"request_data" `
	RequestUrl      string                 `json:"request_url"  mapstructure:"request_url"`
	RequestDataUrl  string                 `json:"request_data_url"  mapstructure:"request_data_url"`
	RequestMethod   string                 `json:"type"  mapstructure:"type"`
	TokenKey        string                 `json:"addr"  mapstructure:"addr"`
	Mark            string                 `json:"mark"  mapstructure:"mark"`
	C               *check.C
	J               *json_key.J
	common.BaseJccAPITesting
}

func (jt *JsonTesting) UrlPath() (s string) { return jt.RequestUrl }

/**
 * Post参数处理
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) HandleParam(gat common.GoApiTesting) map[string]interface{} {
	// 与全局参数合并
	for k, v := range jt.J.Params {
		if _, ok := jt.RequestData[k]; !ok {
			jt.RequestData[k] = v
		}
	}

	common.ParamRandomByMap(&jt.RequestData)
	return jt.RequestData
}

/**
 * Post参数处理
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) HandleHeader(gat common.GoApiTesting) map[string]string {
	// 与全局参数合并
	for k, v := range jt.J.Header {
		if _, ok := jt.Header[k]; !ok {
			jt.Header[k] = v
		}
	}

	return jt.Header
}

/**
 * Get参数处理
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) HandleUrlCode(gat common.GoApiTesting) map[string]interface{} {
	// 与全局参数合并
	for k, v := range jt.J.Params {
		if _, ok := jt.RequestData[k]; !ok {
			jt.RequestData[k] = v
		}
	}

	return jt.RequestData
}

func (jt *JsonTesting) SetUpTest(c *check.C) {
	// fmt.Println("-----------SetUpTest----" + fmt.Sprint() + "---------------")
	common.SetToken(jt.J)
	jt.Validation()
}
func (jt *JsonTesting) TearDownTest(c *check.C) {
	// fmt.Println("-----------TearDownTest----" + fmt.Sprint() + "---------------")
	jt.Validation()
}

func (jt *JsonTesting) TearDownSuite(c *check.C) {
}
func (jt *JsonTesting) SetUpSuite(c *check.C) {
	// if !isLogin {
	// 	jt.C = c
	// 	// fmt.Println("-----------JsonTesting-SetUp---" + fmt.Sprint() + "---------------")
	// 	login := new(Login)
	// 	login.TestLoginSuccess(jt.C)
	// 	pmToken := login.Response.JSON().Object().Raw()["token"].(string)
	// 	jt.AddHeader("PmToken", pmToken)
	// 	jt.AddParams("PmToken", pmToken)
	//
	// 	storeLogin := new(StoreLogin)
	// 	storeLogin.TestLoginSuccess(jt.C)
	// 	token := storeLogin.Response.JSON().Object().Raw()["token"].(string)
	// 	jt.AddHeader("Token", token)
	// 	jt.AddParams("Token", token)
	//
	// 	// jt.AddHeader("PmToken", "e611c721-a4d2-40cb-a8c3-6c8995e495b5")
	// 	// jt.AddHeader("Token", "a2ae8a7e-b705-43ea-bd38-c91721aab743")
	// 	// jt.AddHeader("MemberToken", "a2ae8a7e-b705-43ea-bd38-c91721aab743")
	// 	isLogin = true
	// }

}

func (jt *JsonTesting) SetUp() {
	// fmt.Println("-----------SetUp----" + fmt.Sprint() + "---------------")
}
func (jt *JsonTesting) TearDown() {
	// fmt.Println("-----------TearDown----" + fmt.Sprint() + "---------------")
}

/**
 * 获取路由地址
 * @Author: cs_shuai
 * @Date: 2020-08-11
 */
func (jt *JsonTesting) GetRouteDir() string {
	return viper.GetString("JSON_ROUTE_PATH")
}

/**
 * 初始化测试
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func (jt *JsonTesting) NewTesting(requestDatas []interface{}) common.AutoTesting {
	// 处理并注册到测试
	for _, requestData := range requestDatas {
		// 初始化
		jt.Initialization()
		newJsonTesting := new(JsonTesting)
		if err := mapstructure.Decode(requestData, newJsonTesting); err != nil {
			panic(err)
		}

		newJsonTesting.RequestDataList = append(newJsonTesting.RequestDataList, newJsonTesting.RequestData)
		// 存在请求参数地址 获取参数并合并
		if newJsonTesting.RequestDataUrl != "" {
			paramArr := common.GetParamsByJsonFile(newJsonTesting.RequestDataUrl, viper.GetString("JSON_PATH"))
			newJsonTesting.RequestDataList = append(newJsonTesting.RequestDataList, paramArr...)
		}

		// 注册到测试中
		for _, test := range newJsonTesting.RequestDataList {
			jt.J = json_key.NewJ()
			jTemp := *newJsonTesting
			jTemp.RequestData = test.(map[string]interface{})
			jTemp.J = jt.J
			jTemp.J.Params = jTemp.RequestData
			jt.J.ParamsList[jt.RequestUrl] = append(jt.J.ParamsList[jt.RequestUrl], jt.RequestData)

			check.Suite(&jTemp)
		}
	}

	return jt
}

/**
 * 测试执行方法 [通过json生成不同数据的对象 最后执行的方法]
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func (jt *JsonTesting) TestRun(c *check.C) common.AutoTesting {
	// fmt.Println("-----TestRun----------" + fmt.Sprint(*jt) + "---------------")
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
 * 调取验证类
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func (jt *JsonTesting) Validation() {
	funcName := getFuncName(2)

	for _, v := range json_key.JsonKeyExpand {
		switch funcName {
		case "SetUpTest":
			jt.MakeMark()
			jt.Header = jt.J.Header
			jt.J.HeaderList[jt.RequestUrl] = append(jt.J.HeaderList[jt.RequestUrl], jt.J.Header)
			jt.J.HeaderMap[jt.Mark] = jt.J.Header

			jt.J.Params = jt.RequestData
			jt.J.ParamsList[jt.RequestUrl] = append(jt.J.ParamsList[jt.RequestUrl], jt.RequestData)
			jt.J.ParamsMap[jt.Mark] = jt.RequestData
			if value, ok := jt.RequestData[v.GetJsonKey()]; ok {
				v.SetJsonValue(value)
				v.SetUpRun(jt.J)
			}
		case "TearDownTest":
			jt.J.Response = jt.Response
			jt.J.ResponseList[jt.RequestUrl] = append(jt.J.ResponseList[jt.RequestUrl], jt.Response)
			jt.J.ResponseMap[jt.Mark] = jt.Response
			if value, ok := jt.RequestData[v.GetJsonKey()]; ok {
				v.SetJsonValue(value)
				v.TearDownRun(jt.J)
			}
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

func (jt *JsonTesting) MakeMark() {
	if jt.Mark == "" {
		str := fmt.Sprint(jt.Header) + fmt.Sprint(jt.RequestData) + jt.RequestUrl
		h := sha1.New()
		h.Write([]byte(str))
		jt.Mark = fmt.Sprintf("%x", h.Sum(nil))
	}

}
