package validation

import (
	"github.com/gavv/httpexpect"
)

/**
 * 返回匹配
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type After struct {
	Key   string
	Value map[string]interface{}
}

func (a *After) GetJsonKey() string {
	a.Key = "after"
	return a.Key
}

func (a *After) SetJsonValue(value interface{}) {
	a.Value = value.(map[string]interface{})
}

func (a *After) GetJsonValue() interface{} {
	return a.Value
}

func (a *After) SetUpRun(params *map[string]interface{}) {

}

func (a *After) TearDownRun(res *httpexpect.Response, params *map[string]interface{}) {}
