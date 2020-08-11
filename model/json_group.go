package model

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
	"jccAPITest/common"
	"sync"
)

const After = "after"
const Before = "before"

/**
 * 分组自动测试
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type JsonGroupTesting struct {
	G           *g
	JsonTesting `mapstructure:",squash"`
}

type g struct {
	IsMain     bool
	Header     map[string]string
	Params     map[string]interface{}
	Wait       sync.WaitGroup
	ParamsChan chan map[string]interface{}
}

func NewG() *g {
	g := new(g)
	g.Header = make(map[string]string)
	g.Params = make(map[string]interface{})
	var wg sync.WaitGroup
	g.Wait = wg
	g.ParamsChan = make(chan map[string]interface{})

	return g
}

func (jgt *JsonGroupTesting) AddHeader(key, value string) {
	jgt.G.Header[key] = value
}

func (jgt *JsonGroupTesting) AddParams(key string, value interface{}) {
	jgt.G.Params[key] = value
}

func (jgt *JsonGroupTesting) GetHeader() map[string]string {
	return jgt.G.Header
}

func (jgt *JsonGroupTesting) GetParams() map[string]interface{} {
	return jgt.G.Params
}

func (jgt *JsonGroupTesting) NewTesting(requestData interface{}) common.AutoTesting {
	newJsonTesting := new(JsonGroupTesting)
	conf := new(check.RunConf)
	jgt.G = NewG()
	if err := mapstructure.WeakDecode(requestData, newJsonTesting); err != nil {
		panic(err)
	}
	// 存在请求参数地址 获取参数并合并
	if newJsonTesting.RequestDataUrl != "" {
		paramArr := common.GetParamsByJsonFile(newJsonTesting.RequestDataUrl, viper.GetString("JSON_PATH"))
		newJsonTesting.RequestDataList = append(newJsonTesting.RequestDataList, paramArr...)
	}

	// 开启信道
	jgt.ParamChan()

	// 注册到测试中
	for key, test := range newJsonTesting.RequestDataList {
		jTemp := *newJsonTesting
		jTemp.RequestData = test.(map[string]interface{})
		jTemp.G = jgt.G

		if key == len(newJsonTesting.RequestDataList)-1 {
			// 等待前置任务执行完毕
			jgt.G.Wait.Wait()
			jgt.G.IsMain = true
			// fmt.Println("--------注册到测试-------" + fmt.Sprint(jgt.G.Params) + "---------------")
			check.Suite(&jTemp)
		} else {
			jgt.G.Wait.Add(1)
			check.Run(&jTemp, conf)
		}
	}

	return newJsonTesting
}
func (jgt *JsonGroupTesting) SetUpTest(c *check.C) {
	// fmt.Println("------Params---------" + fmt.Sprint(jgt.G.Params) + "---------------")
	// 从前置参数中获取 Before 的值
	for key, value := range jgt.RequestData {
		if value == Before {
			if beforeValue, ok := jgt.G.Params[key]; ok {
				jgt.RequestData[key] = beforeValue
			}
		}
	}
	// fmt.Println("----Before--RequestData---------" + fmt.Sprint(jgt.RequestData) + "---------------")

	jgt.Validation()
}

func (jgt *JsonGroupTesting) TearDownTest(c *check.C) {
	if !jgt.G.IsMain {
		jgt.G.Wait.Done()
	}
	// 处理前置参数
	jgt.HandleAfter()
	// 建议处理
	jgt.Validation()
}

/**
 * 参数信道
 * @Author: cs_shuai
 * @Date: 2020-08-11
 */
func (jgt *JsonGroupTesting) ParamChan() {
	go func() {
		for {
			select {
			case m := <-jgt.G.ParamsChan:
				// fmt.Println("-----ParamsChan----------" + fmt.Sprint(m) + "---------------")
				jgt.AddParamsGloda(m)
			}
		}
	}()
}

/**
 * 处理前置参数
 * @Author: cs_shuai
 * @Date: 2020-08-11
 */
func (jgt *JsonGroupTesting) HandleAfter() {
	if _, ok := jgt.RequestData[After]; ok {
		jgt.G.ParamsChan <- jgt.RequestData[After].(map[string]interface{})
	}
}

func (jgt *JsonGroupTesting) AddParamsGloda(m map[string]interface{}) {
	for key, value := range m {
		jgt.G.Params[key] = value
	}
}
