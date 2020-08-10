package validation

import "github.com/gavv/httpexpect"

var CheckList []BaseInterface

const (
	SetUpSuite    = "SetUpSuite"
	TearDownSuite = "TearDownSuite"
	SetUpTest     = "SetUpTest"
	TearDownTest  = "TearDownTest"
)

/**
 * 验证接口
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type BaseInterface interface {
	GetJsonKey() string
	GetRunFunc() string
	SetJsonValue(interface{})
	GetJsonValue() interface{}
	Run(res *httpexpect.Response, params *map[string]interface{})
}

/**
 * 注册
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func Register(v ...BaseInterface) {
	CheckList = append(CheckList, v...)
}
