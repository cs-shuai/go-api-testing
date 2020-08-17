package json_value

import (
	"math/rand"
	"time"
)

func init() {
	RegisterJsonValue(new(Auto))
}

type Auto struct {
}

func (a *Auto) GetJsonValue() string {
	return "auto"
}

func (a *Auto) Run() interface{} {
	return randomString(8)
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
