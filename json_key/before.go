package json_key

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"math/rand"
	"strings"
)

func init() {
	RegisterJsonKey(new(Before))
}

/**
 * 返回匹配
 * @Author: cs_shuai
 * @Date: 2020-08-10
 */
type Before struct {
	Key   string
	Value []*before
}

type before struct {
	Url       string `json:"url"  mapstructure:"url"`
	Key       string `json:"key"  mapstructure:"key"`
	IsHeader  bool   `json:"is_header"  mapstructure:"is_header"`
	BeforeKey string `json:"before_key"  mapstructure:"before_key"`
}

func (b *Before) GetJsonKey() string {
	b.Key = "before"
	return b.Key
}

func (b *Before) SetJsonValue(value interface{}) {
	// fmt.Println("----Before---value--------" + fmt.Sprint(value) + "---------------")
	var valueList []*before

	for _, v := range value.([]interface{}) {
		b := new(before)
		if err := mapstructure.Decode(v, b); err != nil {
			panic(err)
		}
		valueList = append(valueList, b)
	}

	b.Value = valueList
}

func (b *Before) GetJsonValue() interface{} {
	return b.Value
}

func (b *Before) SetUpRun(params *J) {
	for _, beforeInfo := range b.Value {
		responseInfo := params.ResponseList[beforeInfo.Url][0]
		res := responseInfo.JSON().Raw()
		// fmt.Println("--------res-------" + fmt.Sprint(res) + "---------------")
		keyArr := strings.Split(beforeInfo.BeforeKey, ".")

		v := printMap(res.(map[string]interface{}), keyArr, 0)
		// fmt.Println("--------res-vvvv------" + fmt.Sprint(v) + "---------------")

		if beforeInfo.IsHeader {
			params.Header[beforeInfo.Key] = fmt.Sprint(v)
		} else {
			params.Params[beforeInfo.Key] = v
		}
	}
	// fmt.Println("----Before---params--------" + fmt.Sprint(params) + "---------------")

	delete(params.Params, b.GetJsonKey())
}

func (b *Before) TearDownRun(params *J) {}

/**
 * 格式化map
 * @Author: cs_shuai
 * @Date: 2020-08-15
 */
func printMap(m map[string]interface{}, keys []string, key int) (r interface{}) {
	r = m[keys[key]]
	switch m[keys[key]].(type) {
	case []interface{}:
		var i int
		if keys[key+1] == "random" {
			l := len(m[keys[key]].([]interface{}))
			i = rand.Intn(l - 1)
		}
		r = printMap(m[keys[key]].([]interface{})[i].(map[string]interface{}), keys, key+2)

	case map[string]interface{}:
		r = printMap(m[keys[key]].(map[string]interface{}), keys, key+1)
	}

	return r
}
