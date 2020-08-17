package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cs-shuai/go-api-testing/json_key"
	"github.com/cs-shuai/go-api-testing/json_value"
	"github.com/gavv/httpexpect"
	"github.com/spf13/viper"
	"gopkg.in/check.v1"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

func init() {
	ConfigInit()
}

// 注册校验
func RegisterCheck(tests ...GoApiTesting) {
	for _, test := range tests {
		var _ = check.Suite(test)
		test.Initialization()
		TestList = append(TestList, test)
	}
}

/**
 * GET接口请求
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func HttpGet(c *check.C, t GoApiTesting) *httpexpect.Response {
	t.Initialization()
	uri := urlHandle(t)

	m := t.HandleUrlCode(t)

	lc := NewLoggerReporter(c, m)
	e := httpexpect.New(lc, Host)

	// fmt.Println("----------t.UrlPath() -----" + fmt.Sprint(uri) + "---------------")
	// fmt.Println("----------t.Token() -----" + fmt.Sprint(Token) + "---------------")

	request := e.GET(uri)

	// 参数处理
	for key, value := range m {
		request.WithQuery(key, value)
	}

	// 头处理
	for key, value := range t.HandleHeader(t) {
		// fmt.Println("------------WithHeader---" + fmt.Sprint(key) + "---------------")
		// fmt.Println("------------WithHeader---" + fmt.Sprint(value) + "---------------")
		request.WithHeader(key, value)
	}
	r := request.Expect().
		Status(http.StatusOK)
	t.SetResponse(r)

	return r
}

/**
 * POST接口请求
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func HttpPost(c *check.C, t GoApiTesting) *httpexpect.Response {
	t.Initialization()

	contentType := "application/x-www-form-urlencoded;charset=utf-8"
	r := httpPost(c, t, contentType)

	t.SetResponse(r)

	return r
}

/**
 * post请求
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func httpPost(c *check.C, t GoApiTesting, contentType string) *httpexpect.Response {
	// 域名
	if Host == "" {
		panic("host is null")
	}

	// 请求地址处理
	uri := urlHandle(t)
	// 请求参数处理
	m := t.HandleParam(t)

	lc := NewLoggerReporter(c, m)

	e := httpexpect.New(lc, Host)
	// fmt.Println("-----------uri----" + fmt.Sprint(uri) + "---------------")
	// fmt.Println("------------mmm---" + fmt.Sprint(m) + "---------------")
	// fmt.Println("------------Token---" + fmt.Sprint(Token) + "---------------")
	// fmt.Println("------------contentType---" + fmt.Sprint(contentType) + "---------------")
	request := e.POST(uri)

	// 头处理
	for key, value := range t.HandleHeader(t) {
		// fmt.Println("------------WithHeader---" + fmt.Sprint(key) + "---------------")
		// fmt.Println("------------WithHeader---" + fmt.Sprint(value) + "---------------")

		request.WithHeader(key, value)
	}

	// fmt.Println("-----------请求参数----" + fmt.Sprint(m) + "---------------")
	r := request.
		WithHeader("ContentType", contentType). // 定义头信息
		WithForm(m).
		Expect().
		Status(http.StatusOK)
	return r
}

/**
 * 请求地址处理
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func urlHandle(t GoApiTesting) string {
	uri := t.UrlPath()
	if uri == "" {
		panic("uri is null")
	}
	uri = strings.Trim(uri, "/")
	uri = "/" + uri

	return uri
}

/**
 * POST:JSON接口请求
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func HttpPostJson(c *check.C, t GoApiTesting) *httpexpect.Response {
	t.Initialization()

	contentType := "application/json;charset=utf-8"
	r := httpPost(c, t, contentType)

	t.SetResponse(r)

	return r
}

/**
 * 结构体转urlCode
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func structToUrlCode(t GoApiTesting) map[string]interface{} {
	m := make(map[string]interface{})

	sv := reflect.ValueOf(t).Elem()
	st := reflect.TypeOf(t).Elem()
	for i := 0; i < st.NumField(); i++ {
		// fmt.Println("--------key-------" + fmt.Sprint(st.Field(i).Tag.Get("json")) + "---------------")
		// fmt.Println("--------value-------" + fmt.Sprint(sv.Field(i).String()) + "---------------")
		if st.Field(i).Tag.Get("json") != "" {
			m[st.Field(i).Tag.Get("json")] = sv.Field(i).String()
		}
	}

	return m
}

/**
 * 全部随机
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func ParamRandomByMap(m *map[string]interface{}) {
	mapTemp := *m
	res := printInterface(mapTemp, 0)
	mapTemp = res.(map[string]interface{})
	m = &mapTemp
}

/**
 * 格式化数据
 * @Author: cs_shuai
 * @Date: 2020-08-15
 */
func printInterface(i interface{}, skip int) interface{} {
	var r interface{}
	switch i.(type) {
	case []interface{}:
		var temp = make([]interface{}, 0)
		for _, v := range i.([]interface{}) {
			nv := printInterface(v, skip+1)
			temp = append(temp, nv)
		}
		if skip == 1 {
			strJons, _ := json.Marshal(temp)
			r = string(strJons)
		} else {
			r = temp
		}

	case map[string]interface{}:
		for k, v := range i.(map[string]interface{}) {
			temp := printInterface(v, skip+1)
			i.(map[string]interface{})[k] = temp
		}
		if skip == 1 {
			strJons, _ := json.Marshal(i)
			r = string(strJons)
		} else {
			r = i
		}
	default:
		if ev, ok := json_value.JsonValueExpand[fmt.Sprint(i)]; ok {
			if skip == 1 {
				r = fmt.Sprint(ev.Run())
			} else {
				r = ev.Run()
			}
		} else {
			if skip == 1 {
				r = fmt.Sprint(i)
			} else {
				r = i
			}
		}
	}

	return r
}

/**
 * 从json文件中获取参数
 * @Author: cs_shuai
 * @Date: 2020-08-05
 */
func ParamByJson(t GoApiTesting, filename string) {
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
 * 参数随机
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func paramRandom(t GoApiTesting) {
	sv := reflect.ValueOf(t).Elem()
	st := reflect.TypeOf(t).Elem()
	for i := 0; i < st.NumField(); i++ {
		// fmt.Println("--------key-------" + fmt.Sprint(st.Field(i).Tag.Get("json")) + "---------------")
		// fmt.Println("--------value-------" + fmt.Sprint(sv.Field(i).String()) + "---------------")
		if sv.Field(i).String() == "auto" {
			// fmt.Println("---------------" + fmt.Sprint(st.Field(i).Type.String()) + "---------------")
			switch st.Field(i).Type.String() {
			case "string":
				str := randomString(8)
				// fmt.Println("---------------" + fmt.Sprint(str) + "---------------")
				sv.Field(i).SetString(str)
			case "int":
				_int := rand.Intn(10)
				sv.Field(i).SetInt(int64(_int))
			case "int64":
				_int := rand.Intn(10)
				sv.Field(i).SetInt(int64(_int))
			}
		}
	}
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

/**
 * 处理文件名
 * @Author: cs_shuai
 * @Date: 2020-08-06
 */
func HandleFileName(file string) (string, string) {
	fileNameAll := path.Base(file)
	fileSuffix := path.Ext(file)
	filePrefix := fileNameAll[0 : len(fileNameAll)-len(fileSuffix)]

	return filePrefix, fileSuffix
}

/**
 * 检查json文件类型转换为结构体
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func checkFileTypeToStruct(fileName, path string) *jsonFile {
	filePrefix, fileSuffix := HandleFileName(fileName)
	jf := new(jsonFile)
	jf.Path = path + fileName
	jf.FileName = fileName
	jf.FilePrefix = filePrefix
	jf.FileSuffix = fileSuffix

	return jf
}

/**
 * 通过Json获取测试参数
 * @Author: cs_shuai
 * @Date: 2020-08-07
 */
func GetParamsByJsonFile(fileName, path string) (result []interface{}) {
	jf := checkFileTypeToStruct(fileName, path)

	filename := RootPath + path + jf.FileName
	content := ReadJson(filename)
	// fmt.Println(string(content))
	err := json.Unmarshal([]byte(content), &result)
	if err != nil {
		panic(err)
	}

	return result
}

/**
 * 获取文件数据
 * @Author: cs_shuai
 * @Date: 2020-08-15
 */
func GetParamsByJsonFileStruct(jf *jsonFile) (result []interface{}) {
	filename := jf.Path
	content := ReadJson(filename)
	// fmt.Println(string(content))
	err := json.Unmarshal([]byte(content), &result)
	if err != nil {
		panic(err)
	}

	return result
}

/**
 * 读取Json
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func ReadJson(filePath string) (result string) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	buf := bufio.NewReader(file)
	for {
		s, err := buf.ReadString('\n')
		if strings.HasPrefix(s, "//") {
			continue
		}
		result += s
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
	}
	return result
}

/**
 * 获取全部Json
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func ReadAllJson(filename string) string {
	fileObj, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fileObj.Close()
	content, err := ioutil.ReadAll(fileObj)
	return string(content)
}

/**
 * 设置Token
 * @Author: cs_shuai
 * @Date: 2020-08-15
 */
func SetToken(j *json_key.J) {
	tokenArr := viper.Get("token")
	headTokenArr := viper.GetStringMapString("head_token")
	for k, v := range tokenArr.(map[string]interface{}) {
		j.Params[k] = v
	}
	for k, v := range headTokenArr {
		j.Header[k] = v
	}
}
