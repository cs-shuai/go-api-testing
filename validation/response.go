package validation

import (
	"github.com/gavv/httpexpect"
	"jccAPITest/common"
)

func init() {
	Register(new(Response))
}

/**
 * 返回匹配
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type Response struct {
	Key   string
	Value interface{}
}

func (r *Response) GetJsonKey() string {
	r.Key = "response"
	return r.Key
}

func (r *Response) GetRunFunc() string {
	return TearDownTest
}

func (r *Response) SetJsonValue(value interface{}) {
	r.Value = value
}

func (r *Response) GetJsonValue() interface{} {
	return r.Value
}

func (r *Response) Run(res *httpexpect.Response, params *map[string]interface{}) {
	responseKey := common.ResponseKey
	response := common.Response
	equalValue := "成功"
	var tempMap = *params
	if value, ok := tempMap[response]; ok {
		equalValue = value.(string)
	}
	params = &tempMap
	res.JSON().Object().Value(responseKey).Equal(equalValue)
}
