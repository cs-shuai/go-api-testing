package model

import (
	"github.com/cs-shuai/go-api-testing/common"
	"github.com/cs-shuai/go-api-testing/json_key"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
)

/**
 * 分组自动测试
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type JsonGroupTesting struct {
	JsonTesting `mapstructure:",squash"`
}

func (jgt *JsonGroupTesting) NewTesting(requestDatas []interface{}) common.AutoTesting {
	jgt.J = json_key.NewJ()
	// 处理单个文件下的多数据
	for key, requestData := range requestDatas {
		newJsonTesting := new(JsonGroupTesting)
		conf := new(check.RunConf)
		if err := mapstructure.WeakDecode(requestData, newJsonTesting); err != nil {
			panic(err)
		}

		if newJsonTesting.RequestUrl == "" {
			continue
		}

		// 存在请求参数地址 获取参数并合并
		if newJsonTesting.RequestDataUrl != "" {
			paramArr := common.GetParamsByJsonFile(newJsonTesting.RequestDataUrl, viper.GetString("JSON_PATH"))
			newJsonTesting.RequestDataList = append(newJsonTesting.RequestDataList, paramArr...)
		}

		// 注册到测试中
		for _, test := range newJsonTesting.RequestDataList {
			jTemp := *newJsonTesting
			jTemp.RequestData = test.(map[string]interface{})
			jTemp.J = jgt.J
			if key == len(requestDatas)-1 {
				check.Suite(&jTemp)
			} else {
				check.Run(&jTemp, conf)
			}
		}
	}
	return jgt
}
func (jgt *JsonGroupTesting) SetUpTest(c *check.C) {
	jgt.Validation()
}

func (jgt *JsonGroupTesting) TearDownTest(c *check.C) {
	jgt.Validation()
}

func (jgt *JsonGroupTesting) AddParamsGloda(m map[string]interface{}) {
	for key, value := range m {
		jgt.J.Params[key] = value
	}
}

/**
 * 获取路由地址
 * @Author: cs_shuai
 * @Date: 2020-08-11
 */
func (jgt *JsonGroupTesting) GetRouteDir() string {
	return viper.GetString("JSON_GROUP_ROUTE_PATH")
}
