package json_key

import "github.com/spf13/viper"

var ResponseMessage string

func init() {
	RegisterJsonKey(new(Response))
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
func (r *Response) SetJsonValue(value interface{}) {
	r.Value = value
}

func (r *Response) GetJsonValue() interface{} {
	return r.Value
}

func (r *Response) TearDownRun(params *J) {
	if ResponseMessage == "" {
		ResponseMessage = viper.GetString("RESPONSE_MESSAGE")
	}
	params.Response.JSON().Object().Value(ResponseMessage).Equal(r.Value)
}

func (r *Response) SetUpRun(params *J) {

}
