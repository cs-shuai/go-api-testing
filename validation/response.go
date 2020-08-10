package validation

import (
	"fmt"
	"github.com/gavv/httpexpect"
	"jccAPITest/common"
)

func init() {
	Register(new(Response))
}

type Response struct {
	Key   string
	Value string
}

func (r *Response) GetJsonKey() string {
	r.Key = "response"
	return r.Key
}

func (r *Response) GetRunFunc() string {
	return TearDownTest
}

func (r *Response) SetJsonValue(value interface{}) {
	r.Value = fmt.Sprint(value)
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
