package json_value

import (
	"math/rand"
)

func init() {
	RegisterJsonValue(new(AutoInt))
}

type AutoInt struct {
}

func (a *AutoInt) GetJsonValue() string {
	return "auto_int"
}

func (a *AutoInt) Run() interface{} {
	return rand.Intn(10)
}
