package model

type JsonGroupTesting struct {
	Header map[string]string
	Params map[string]interface{}
	JsonTesting
}

func (jgt *JsonGroupTesting) AddHeader(key, value string) {
	jgt.Header[key] = value
}

func (jgt *JsonGroupTesting) AddParams(key string, value interface{}) {
	jgt.Params[key] = value
}

func (jgt *JsonGroupTesting) GetHeader() map[string]string {
	return jgt.Header
}

func (jgt *JsonGroupTesting) GetParams() map[string]interface{} {
	return jgt.Params
}
