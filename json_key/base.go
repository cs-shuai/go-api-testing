package json_key

import (
	"github.com/gavv/httpexpect"
)

var JsonKeyExpand []BaseJsonKeyExpandInterface

/**
 * 验证接口
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type BaseJsonKeyExpandInterface interface {
	GetJsonKey() string
	SetJsonValue(interface{})
	GetJsonValue() interface{}
	SetUpRun(*J)
	TearDownRun(*J)
}

/**
 * 注册
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func RegisterJsonKey(v ...BaseJsonKeyExpandInterface) {
	JsonKeyExpand = append(JsonKeyExpand, v...)
}

type J struct {
	IsMain       bool
	Header       map[string]string
	Params       map[string]interface{}
	Response     *httpexpect.Response
	ResponseMap  map[string]*httpexpect.Response
	ResponseList map[string][]*httpexpect.Response
	ParamsMap    map[string]map[string]interface{}
	ParamsList   map[string][]map[string]interface{}
	HeaderMap    map[string]map[string]string
	HeaderList   map[string][]map[string]string
}

/**
 * 初始化Josn通用对象
 * @Author: cs_shuai
 * @Date: 2020-08-15
 */
func NewJ() *J {
	j := new(J)
	j.Header = make(map[string]string)
	j.Params = make(map[string]interface{})
	j.ResponseList = make(map[string][]*httpexpect.Response)
	j.ParamsList = make(map[string][]map[string]interface{})
	j.HeaderList = make(map[string][]map[string]string)
	j.HeaderMap = make(map[string]map[string]string)
	j.ParamsMap = make(map[string]map[string]interface{})
	j.ParamsMap = make(map[string]map[string]interface{})
	j.ResponseMap = make(map[string]*httpexpect.Response)

	return j
}
