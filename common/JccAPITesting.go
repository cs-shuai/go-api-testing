package common

import (
	"encoding/json"
	"fmt"
	"github.com/gavv/httpexpect"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"
)

var URL string
var TokenKey string
var Token string
var TestList []JccAPITesting

func init() {
	// 配置初始化
	configInit()
}

/**
 * 配置初始化
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func configInit() {
	v := viper.New()
	v.SetConfigName("config") // name of config file (without extension)
	v.SetConfigType("toml")   // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath("../conf/")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("no such config file")
		} else {
			// Config file was found but another error was produced
			log.Println("read config error")
		}
		log.Fatal(err) // 读取配置文件失败致命错误
	}

	URL = v.GetString("URL")
	TokenKey = v.GetString("TOKEN_KEY")
}

// 集餐厨接口测试接口类
type JccAPITesting interface {
	Initialization()
	SetUpSuite(c *check.C)
	TearDownSuite(c *check.C)
	SetUpTest(c *check.C)
	TearDownTest(c *check.C)
	UrlPath() string
	SetResponse(*httpexpect.Response)
}

// 测试基础类
type BaseJccAPITesting struct {
	Token    string               `json:"token"`
	Response *httpexpect.Response `json:"-"`
}

func (t *BaseJccAPITesting) Initialization() {
}

func (BaseJccAPITesting) SetUpSuite(c *check.C) {
}

func (BaseJccAPITesting) TearDownSuite(c *check.C) {
}

func (BaseJccAPITesting) SetUpTest(c *check.C) {
}

func (t *BaseJccAPITesting) SetResponse(response *httpexpect.Response) {
	t.Response = response
}

func (BaseJccAPITesting) TearDownTest(c *check.C) {
}

func (t *BaseJccAPITesting) UrlPath() string {
	return ""
}

// 注册校验
func RegisterCheck(tests ...JccAPITesting) {
	for _, test := range tests {
		var _ = check.Suite(test)
		fmt.Println(test.UrlPath())
		TestList = append(TestList, test)
	}
}

/**
 * GET接口请求
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func HttpGet(c *check.C, t JccAPITesting) *httpexpect.Response {
	t.Initialization()
	paramsStr := structToUrlCode(t)
	fmt.Println("----------paramsStr-----" + fmt.Sprint(paramsStr) + "---------------")
	fmt.Println("----------t.UrlPath() -----" + fmt.Sprint(t.UrlPath()) + "---------------")
	fmt.Println("----------t.Token() -----" + fmt.Sprint(Token) + "---------------")
	fmt.Println("----------t.TokenKey() -----" + fmt.Sprint(TokenKey) + "---------------")
	uri := t.UrlPath() + "?" + paramsStr
	e := httpexpect.New(c, URL)
	r := e.GET(uri).
		Expect().
		Status(http.StatusOK)
	t.SetResponse(r)

	return r
}

/**
 * POST接口请求
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func HttpPost(c *check.C, t JccAPITesting) *httpexpect.Response {
	t.Initialization()
	var err error
	uri := t.UrlPath()
	e := httpexpect.New(c, URL)
	m := make(map[string]interface{})
	j, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(j, &m)
	if err != nil {
		panic(err)
	}

	contentType := "application/x-www-form-urlencoded;charset=utf-8"
	r := e.POST(uri). // post 请求
				WithHeader(TokenKey, Token).            // 定义头信息
				WithHeader("ContentType", contentType). // 定义头信息
				WithForm(m).
				Expect().
				Status(http.StatusOK)
	t.SetResponse(r)

	return r
}

/**
 * POST:JSON接口请求
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func HttpPostJson(c *check.C, t JccAPITesting) *httpexpect.Response {
	t.Initialization()

	var err error

	uri := t.UrlPath()
	e := httpexpect.New(c, URL)
	contentType := "application/json;charset=utf-8"
	m := make(map[string]interface{})
	j, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(j, &m)
	if err != nil {
		panic(err)
	}

	r := e.POST(uri). // post 请求
				WithHeader(TokenKey, Token).            // 定义头信息
				WithHeader("ContentType", contentType). // 定义头信息
				WithJSON(m).                            // 传入json body
				Expect().
				Status(http.StatusOK)
	t.SetResponse(r)

	return r
}

/**
 * 结构体转urlCode
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func structToUrlCode(t JccAPITesting) string {
	sv := reflect.ValueOf(t).Elem()
	st := reflect.TypeOf(t).Elem()
	params := url.Values{}
	for i := 0; i < st.NumField(); i++ {
		// fmt.Println("--------key-------" + fmt.Sprint(st.Field(i).Tag.Get("json")) + "---------------")
		// fmt.Println("--------value-------" + fmt.Sprint(sv.Field(i).String()) + "---------------")
		if st.Field(i).Tag.Get("json") != "" {
			params.Add(st.Field(i).Tag.Get("json"), sv.Field(i).String())
		}
	}

	return params.Encode()
}

/**
 * 设置TOKEN
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func SetHeaderToken(token string) {
	fmt.Println("---------token------" + fmt.Sprint(token) + "---------------")
	Token = token
}

/**
 * 从json文件中获取参数
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func ParamByJson(t JccAPITesting, filename string) {
	fileObj, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fileObj.Close()
	content, err := ioutil.ReadAll(fileObj)
	fmt.Println(string(content))
	err = json.Unmarshal(content, t)
	// 随机
	paramRandom(t)
}

/**
 * 全部随机
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func paramRandomByStruct(t JccAPITesting) {
	sv := reflect.ValueOf(t).Elem()
	st := reflect.TypeOf(t).Elem()
	for i := 0; i < st.NumField(); i++ {
		// fmt.Println("--------key-------" + fmt.Sprint(st.Field(i).Tag.Get("json")) + "---------------")
		// fmt.Println("--------value-------" + fmt.Sprint(sv.Field(i).String()) + "---------------")
		fmt.Println("---------------" + fmt.Sprint(st.Field(i).Type.String()) + "---------------")
		switch st.Field(i).Type.String() {
		case "string":
			str := randomString(8)
			fmt.Println("---------------" + fmt.Sprint(str) + "---------------")
			sv.Field(i).SetString(str)
		case "int":
			int := rand.Intn(10)
			sv.Field(i).SetInt(int64(int))
		case "int64":
			int := rand.Intn(10)
			sv.Field(i).SetInt(int64(int))
		}
	}

}

func paramRandom(t JccAPITesting) {
	sv := reflect.ValueOf(t).Elem()
	st := reflect.TypeOf(t).Elem()
	for i := 0; i < st.NumField(); i++ {
		// fmt.Println("--------key-------" + fmt.Sprint(st.Field(i).Tag.Get("json")) + "---------------")
		// fmt.Println("--------value-------" + fmt.Sprint(sv.Field(i).String()) + "---------------")
		if sv.Field(i).String() == "auto" {
			fmt.Println("---------------" + fmt.Sprint(st.Field(i).Type.String()) + "---------------")
			switch st.Field(i).Type.String() {
			case "string":
				str := randomString(8)
				fmt.Println("---------------" + fmt.Sprint(str) + "---------------")
				sv.Field(i).SetString(str)
			case "int":
				int := rand.Intn(10)
				sv.Field(i).SetInt(int64(int))
			case "int64":
				int := rand.Intn(10)
				sv.Field(i).SetInt(int64(int))
			}
		}
	}
}

func randomString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
