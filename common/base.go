package common

import (
	"database/sql"
	"encoding/json"
	"github.com/cs-shuai/go-api-testing/json_key"
	"github.com/gavv/httpexpect"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var Host string
var TestList []GoApiTesting
var RootPath string

// 自动测试接口
type AutoTesting interface {
	NewTesting([]interface{}) AutoTesting
	TestRun(*check.C) AutoTesting
	Request() AutoTesting
	SetUp()
	TearDown()
	GetRouteDir() string
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
	HandleHeader(GoApiTesting) map[string]string
	HandleUrlCode(GoApiTesting) map[string]interface{}
	// AddHeader(key, value string)
	// AddParams(key string, value interface{})
	// GetHeader() map[string]string
	// GetParams() map[string]interface{}
}

// 测试基础类
type BaseJccAPITesting struct {
	Token    string               `json:"token"`
	Response *httpexpect.Response `json:"-"`
	Db       *sql.DB
	J        *json_key.J
}

type jsonFile struct {
	Path       string
	FilePrefix string
	FileSuffix string
	FileName   string
	File       os.FileInfo
}

func (t *BaseJccAPITesting) Initialization() {
	db, err := sql.Open("mysql", viper.GetString("SQLCONN"))
	if err != nil {
		panic(err)
	}

	t.Db = db
}

func (t *BaseJccAPITesting) SetUpSuite(_ *check.C) {}

func (t *BaseJccAPITesting) TearDownSuite(_ *check.C) {}

func (t *BaseJccAPITesting) SetUpTest(_ *check.C) {}

func (t *BaseJccAPITesting) SetResponse(response *httpexpect.Response) {
	t.Response = response
}

func (t *BaseJccAPITesting) TearDownTest(_ *check.C) {}

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

func (t *BaseJccAPITesting) HandleHeader(gat GoApiTesting) map[string]string {
	var m = make(map[string]string)
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

// /**
//  * 添加头数据
//  * @Author: cs_shuai
//  * @Date: 2020-08-10
//  */
// func (t *BaseJccAPITesting) AddHeader( key string, value string) {
// 	t.J.Header[key] = value
// }
//
// /**
//  * 添加参数数据
//  * @Author: cs_shuai
//  * @Date: 2020-08-10
//  */
// func (t *BaseJccAPITesting) AddParams(key string, value interface{}) {
// 	t.J.Params[key] = value
// }
//
// /**
//  * 获取头数据
//  * @Author: cs_shuai
//  * @Date: 2020-08-10
//  */
// func (t *BaseJccAPITesting) GetHeader() map[string]string {
// 	return t.J.Header
// }
//
// /**
//  * 获取参数数据
//  * @Author: cs_shuai
//  * @Date: 2020-08-10
//  */
// func (t *BaseJccAPITesting) GetParams() map[string]interface{} {
// 	return t.J.Params
// }

/**
 * 测试输出
 * @Author: cs_shuai
 * @Date: 2020-08-15
 */
type loggerReporter struct {
	C          *check.C
	ParamsJson string
}

/**
 * 初始化
 * @Author: cs_shuai
 * @Date: 2020-08-15
 */
func NewLoggerReporter(c *check.C, m map[string]interface{}) *loggerReporter {
	lc := new(loggerReporter)
	lc.C = c
	b, _ := json.Marshal(m)
	lc.ParamsJson = string(b)
	return lc
}

func (l *loggerReporter) Logf(fmt string, args ...interface{}) {
	l.C.Logf(fmt, args...)
}

func (l *loggerReporter) Errorf(message string, args ...interface{}) {
	l.C.Logf("请求参数: %s", l.ParamsJson)
	l.C.Errorf(message, args...)
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
	Host = viper.GetString("HOST")
	// 根目录
	RootPath, _ = os.Getwd()
	RootPath += "/"
}

/**
 * 获取全部文件包括子目录下的文件
 * @Author: cs_shuai
 * @Date: 2020-08-15
 */
func getAllFile(path string) (fileList []*jsonFile) {
	// fmt.Println("---------------" + fmt.Sprint(path) + "---------------")
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		var resList []*jsonFile
		if f.IsDir() {
			resList = getAllFile(path + "/" + f.Name() + "/")
		} else {
			jf := checkFileTypeToStruct(f.Name(), path)
			resList = append(resList, jf)
		}
		fileList = append(fileList, resList...)
	}

	return
}

/**
 * 自动脚本执行
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func AutoTestRun(testingT *testing.T, ts ...AutoTesting) {
	for _, t := range ts {
		// 获取文件地址
		files := getAllFile(t.GetRouteDir())
		t.SetUp()
		// 读取文件
		for _, f := range files {
			// 获取文件数据
			allArr := GetParamsByJsonFileStruct(f)
			// 初始化测试
			t.NewTesting(allArr)
		}

		// 执行测试
		check.TestingT(testingT)

		t.TearDown()
	}
}
