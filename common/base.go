package common

import (
	"encoding/json"
	"github.com/gavv/httpexpect"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"testing"
)

var Host string
var Token string
var TestList []GoApiTesting
var RootPath string

const (
	Response    = "response"
	ResponseKey = "msg"
)

// 自动测试接口
type AutoTesting interface {
	NewTesting(interface{}) AutoTesting
	TestRun(*check.C) AutoTesting
	Request() AutoTesting
	ResponseCheck() AutoTesting
	SetUp()
	TearDown()
	GetWaitGroup() *sync.WaitGroup
	GoApiTesting
}

// 接口测试接口类
type GoApiTesting interface {
	Initialization()
	SetUpSuite(c *check.C)
	TearDownSuite(c *check.C)
	SetUpTest(c *check.C)
	TearDownTest(c *check.C)
	UrlPath() string
	SetResponse(*httpexpect.Response)
	HandleParam(GoApiTesting) map[string]interface{}
	HandleUrlCode(GoApiTesting) map[string]interface{}
}

// 测试基础类
type BaseJccAPITesting struct {
	Token    string               `json:"token"`
	Response *httpexpect.Response `json:"-"`
}

type jsonFile struct {
	Path       string
	FilePrefix string
	FileSuffix string
	FileName   string
}

func (t *BaseJccAPITesting) Initialization() {}

func (t *BaseJccAPITesting) SetUpSuite(c *check.C) {}

func (t *BaseJccAPITesting) TearDownSuite(c *check.C) {}

func (t *BaseJccAPITesting) SetUpTest(c *check.C) {}

func (t *BaseJccAPITesting) SetResponse(response *httpexpect.Response) {
	t.Response = response
}

func (t *BaseJccAPITesting) TearDownTest(c *check.C) {}

func (t *BaseJccAPITesting) HandleUrlCode(gat GoApiTesting) map[string]interface{} {
	m := make(map[string]interface{})

	sv := reflect.ValueOf(gat).Elem()
	st := reflect.TypeOf(gat).Elem()
	for i := 0; i < st.NumField(); i++ {
		// fmt.Println("--------key-------" + fmt.Sprint(st.Field(i).Tag.Get("json")) + "---------------")
		// fmt.Println("--------value-------" + fmt.Sprint(sv.Field(i).String()) + "---------------")
		if st.Field(i).Tag.Get("json") != "" {
			m[st.Field(i).Tag.Get("json")] = sv.Field(i).String()
		}
	}

	return m
}
func (t *BaseJccAPITesting) HandleParam(gat GoApiTesting) map[string]interface{} {
	var m = make(map[string]interface{})
	j, err := json.Marshal(gat)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(j, &m)
	if err != nil {
		panic(err)
	}

	return m
}

/**
 * 配置初始化
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func ConfigInit() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./conf/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("no such config file")
		} else {
			panic("read config error")
		}
	}

	// 请求地址
	Host = viper.GetString("Host")
	// 根目录
	RootPath, _ = os.Getwd()
	RootPath += "/"
}

/**
 * 自动脚本执行
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func AutoTestRun(testingT *testing.T, t AutoTesting) {
	// 获取文件地址
	files, _ := ioutil.ReadDir(viper.GetString("JSON_ROUTE_PATH"))
	// fmt.Println("---------------" + fmt.Sprint(files) + "---------------")

	// 读取文件
	for _, f := range files {
		// 获取数据
		allArr := GetParamsByJsonFile(f.Name(), viper.GetString("JSON_ROUTE_PATH"))

		// 处理并注册到测试
		for _, requestData := range allArr {
			// fmt.Println("---------------" + fmt.Sprint(requestData) + "---------------")
			t.NewTesting(requestData)
		}
	}

	check.TestingT(testingT)
	t.GetWaitGroup().Wait()
	t.TearDown()
}
