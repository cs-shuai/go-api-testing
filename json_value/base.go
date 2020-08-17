package json_value

var JsonValueExpand = make(map[string]BaseJsonValueExpandInterface)

/**
 * 验证接口
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type BaseJsonValueExpandInterface interface {
	GetJsonValue() string
	Run() interface{}
}

/**
 * 注册
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
func RegisterJsonValue(list ...BaseJsonValueExpandInterface) {
	for _, v := range list {
		JsonValueExpand[v.GetJsonValue()] = v
	}
}
