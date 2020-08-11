package validation

import "github.com/gavv/httpexpect"

var CheckList []BaseValidationInterface

/**
 * 验证接口
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type BaseValidationInterface interface {
	GetJsonKey() string
	SetJsonValue(interface{})
	GetJsonValue() interface{}
	SetUpRun(params *map[string]interface{})
	TearDownRun(res *httpexpect.Response, params *map[string]interface{})
}

/**
 * 注册
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func Register(v ...BaseValidationInterface) {
	CheckList = append(CheckList, v...)
}
